package line

import (
	"bytes"
	"encoding/json"
	"github.com/line-api/model/go/model"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

type TalkService struct {
	client *Client

	conn           *model.FTalkServiceClient
	connCompactMsg *model.FCompactMessageServiceClient
}

func (cl *Client) newTalkService() *TalkService {
	return &TalkService{
		client:         cl,
		conn:           cl.ThriftFactory.newTalkServiceClient(),
		connCompactMsg: cl.ThriftFactory.newCompactMessageServiceClient(),
	}
}

/*
Message functions
*/

type MentionData struct {
	S   string `json:"S"`
	E   string `json:"E"`
	Mid string `json:"M"`
}

func (cl *TalkService) React(msgId string, type_ model.PredefinedReactionType) error {
	id, err := strconv.ParseInt(msgId, 10, 64)
	if err != nil {
		return cl.client.afterError(err)
	}
	req := &model.ReactRequest{
		ReqSeq:       cl.client.RequestSequence,
		MessageId:    id,
		ReactionType: &model.ReactionType{PredefinedReactionType: type_},
	}
	return cl.client.afterError(cl.conn.React(cl.client.ctx, req))
}

func (cl *TalkService) SendText(to, text string) (msg *model.Message, err error) {
	msg, err = cl.conn.SendMessage(cl.client.ctx, cl.client.RequestSequence, &model.Message{
		Text:        text,
		To:          to,
		ContentType: model.ContentType_NONE,
	})
	return msg, cl.client.afterError(err)
}
func GetMidTypeFromID(id string) model.ToType {
	switch {
	case strings.HasPrefix(id, "u"):
		return model.ToType_USER
	case strings.HasPrefix(id, "r"):
		return model.ToType_ROOM
	case strings.HasPrefix(id, "c"):
		return model.ToType_GROUP
	}
	return model.ToType(-1)
}

func (cl *TalkService) SendCompactText(to, text string) (*model.Message, error) {
	msg, err := cl.SendMessageCompact(&model.Message{
		Text:        text,
		ToType:      GetMidTypeFromID(to),
		To:          to,
		From:        cl.client.Profile.Mid,
		ContentType: model.ContentType_NONE,
	})
	return msg, cl.client.afterError(err)
}

func (cl *TalkService) SendMessageCompact(msg *model.Message) (*model.Message, error) {
	newMsg, err := cl.connCompactMsg.SendMessageCompact(cl.client.ctx, cl.client.RequestSequence, msg)
	return newMsg, err
}

func (cl *TalkService) UnsendMessage(id string) error {
	return cl.client.afterError(cl.conn.UnsendMessage(cl.client.ctx, cl.client.RequestSequence, id))
}

func (cl *TalkService) SendMessageWithMention(toID string, msgText string, mids []string) (*model.Message, error) {
	var arr []*MentionData
	mentionText := "@sakura"
	texts := strings.Split(msgText, "@!")
	text := ""
	for i := 0; i < len(mids); i++ {
		text += texts[i]
		arr = append(arr, &MentionData{S: strconv.Itoa(utf8.RuneCountInString(text)), E: strconv.Itoa(utf8.RuneCountInString(text) + 7), Mid: mids[i]})
		text += mentionText
	}
	text += texts[len(texts)-1]
	allData, _ := json.MarshalIndent(arr, "", " ")
	msg := model.NewMessage()
	msg.ContentType = model.ContentType_NONE
	msg.To = toID
	msg.Text = text
	msg.ContentMetadata = map[string]string{"MENTION": "{\"MENTIONEES\":" + string(allData) + "}"}
	ms, err := cl.conn.SendMessage(cl.client.ctx, cl.client.RequestSequence, msg)
	return ms, cl.client.afterError(err)
}
func (cl *TalkService) SendTextMentionByList(to string, msgText string, targets []string) error {
	listMid2 := []string{}
	listChar := msgText + "\n"
	listNum := 0
	loopny := len(targets)/20 + 1
	limiter := 0
	limiter2 := 20
	for a := 0; a < loopny; a++ {
		for c := limiter; c < len(targets); c++ {
			if c < limiter2 {
				listChar += strconv.Itoa(listNum) + ": @!\n"
				listNum = listNum + 1
				listMid2 = append(listMid2, targets[c])
				limiter = limiter + 1
			} else {
				limiter2 = limiter + 20
				break
			}
		}
		_, err := cl.SendMessageWithMention(to, listChar, listMid2)
		if err != nil {
			return cl.client.afterError(err)
		}
		listChar = ""
		listMid2 = []string{}
	}
	return nil
}
func (cl *TalkService) SendContact(toMid, contactMid string) (*model.Message, error) {
	msg := model.NewMessage()
	msg.To = toMid
	msg.ContentType = model.ContentType_CONTACT
	msg.ContentMetadata = map[string]string{"mid": contactMid}
	tmp := "0"
	msg.RelatedMessageId = &tmp
	ms, err := cl.conn.SendMessage(cl.client.ctx, cl.client.RequestSequence, msg)
	return ms, cl.client.afterError(err)
}
func (cl *TalkService) SendChatChecked(groupID, messageID string) error {
	err := cl.conn.SendChatChecked(cl.client.ctx, cl.client.RequestSequence, groupID, messageID, 0)
	return cl.client.afterError(err)
}

type Mentions struct {
	MENTIONEES []MentionData `json:"MENTIONEES"`
}

func ParseMentions(msg *model.Message) []string {
	mentions := Mentions{}
	err := json.Unmarshal([]byte(msg.ContentMetadata["MENTION"]), &mentions)
	if err != nil {
		return []string{}
	}
	var mids []string
	for _, mention := range mentions.MENTIONEES {
		mids = append(mids, mention.Mid)
	}
	return mids
}

/*
Profile functions
*/

func (cl *TalkService) UpdateProfileName(name string) error {
	var req = &model.UpdateProfileAttributesRequest{
		ProfileAttributes: map[model.UpdateProfileAttribute]*model.ProfileContent{
			model.UpdateProfileAttribute_DISPLAY_NAME: {Value: name},
		},
	}
	err := cl.conn.UpdateProfileAttributes(cl.client.ctx, cl.client.RequestSequence, req)
	return cl.client.afterError(err)
}
func (cl *TalkService) UpdateProfileBio(bio string) error {
	var req = &model.UpdateProfileAttributesRequest{
		ProfileAttributes: map[model.UpdateProfileAttribute]*model.ProfileContent{
			model.UpdateProfileAttribute_STATUS_MESSAGE: {Value: bio},
		},
	}
	err := cl.conn.UpdateProfileAttributes(cl.client.ctx, cl.client.RequestSequence, req)
	return cl.client.afterError(err)
}
func (cl *TalkService) GetProfile(reason model.SyncReason) (*model.Profile, error) {
	profile, err := cl.conn.GetProfile(cl.client.ctx, reason)
	if err != nil {
		cl.client.Profile = profile
	}
	return profile, cl.client.afterError(err)
}

func (cl *TalkService) CloneProfile(mid string) error {
	contact, err := cl.GetContact(mid)
	if err != nil {
		return cl.client.afterError(err)
	}
	err = cl.UpdateProfileBio(contact.StatusMessage)
	if err != nil {
		return cl.client.afterError(err)
	}
	err = cl.UpdateProfileName(contact.DisplayName)
	if err != nil {
		return cl.client.afterError(err)
	}
	pPath := cl.client.Profile.Mid + ".jpg"
	err = cl.client.DownloadContactIcon(contact.PicturePath, pPath)
	if err != nil {
		return cl.client.afterError(err)
	}
	err = cl.client.UpdateProfilePicture(pPath)
	if err != nil {
		return cl.client.afterError(err)
	}
	os.Remove(pPath)
	oid, err := cl.client.GetProfileCoverId(mid)
	if err != nil {
		return cl.client.afterError(err)
	}
	err = cl.client.UpdateProfileCoverById(oid)
	if err != nil {
		return cl.client.afterError(err)
	}
	return nil
}

/*
Chat functions
*/

func (cl *TalkService) GetChats(mids []string) ([]*model.Chat, error) {
	req := &model.GetChatsRequest{
		ChatMids:     mids,
		WithMembers:  true,
		WithInvitees: true,
	}
	chats, err := cl.conn.GetChats(cl.client.ctx, req)
	if err != nil {
		return nil, cl.client.afterError(err)
	}
	return chats.Chats, nil
}
func (cl *TalkService) GetChat(mid string) (*model.Chat, error) {
	chats, err := cl.GetChats([]string{mid})
	if err != nil {
		return nil, cl.client.afterError(err)
	}
	if len(chats) > 0 {
		return chats[0], cl.client.afterError(err)
	}
	return nil, err
}

func (cl *TalkService) AcceptChatInvitation(gid string) error {
	req := &model.AcceptChatInvitationRequest{
		ReqSeq:  cl.client.RequestSequence,
		ChatMid: gid,
	}
	_, err := cl.conn.AcceptChatInvitation(cl.client.ctx, req)
	return cl.client.afterError(err)
}
func (cl *TalkService) AcceptChatInvitationAsync(gid string) <-chan error {
	req := &model.AcceptChatInvitationRequest{
		ReqSeq:  cl.client.RequestSequence,
		ChatMid: gid,
	}
	_, err := cl.conn.AcceptChatInvitationAsync(cl.client.ctx, req)
	return err
}
func (cl *TalkService) AcceptChatInvitationByTicket(gid, ticket string) error {
	req := &model.AcceptChatInvitationByTicketRequest{
		ReqSeq:   cl.client.RequestSequence,
		ChatMid:  gid,
		TicketId: ticket,
	}
	_, err := cl.conn.AcceptChatInvitationByTicket(cl.client.ctx, req)
	return cl.client.afterError(err)
}
func (cl *TalkService) AcceptChatInvitationByTicketAsync(gid, ticket string) <-chan error {
	req := &model.AcceptChatInvitationByTicketRequest{
		ReqSeq:   cl.client.RequestSequence,
		ChatMid:  gid,
		TicketId: ticket,
	}
	_, err := cl.conn.AcceptChatInvitationByTicketAsync(cl.client.ctx, req)
	return err
}
func (cl *TalkService) CreateChat(name string, targets []string) (*model.Chat, error) {
	req := &model.CreateChatRequest{
		ReqSeq:         cl.client.RequestSequence,
		ChatType:       model.ChatType_GROUP,
		Name:           name,
		TargetUserMids: SliceToSet(targets),
	}
	chat, err := cl.conn.CreateChat(cl.client.ctx, req)
	if err != nil {
		return nil, cl.client.afterError(err)
	}
	return chat.Chat, nil
}
func (cl *TalkService) InviteIntoChat(gid string, targets []string) error {
	req := &model.InviteIntoChatRequest{
		ReqSeq:         cl.client.RequestSequence,
		ChatMid:        gid,
		TargetUserMids: SliceToSet(targets),
	}
	_, err := cl.conn.InviteIntoChat(cl.client.ctx, req)
	return cl.client.afterError(err)
}
func (cl *TalkService) InviteIntoChatAsync(gid string, targets []string) <-chan error {
	req := &model.InviteIntoChatRequest{
		ReqSeq:         cl.client.RequestSequence,
		ChatMid:        gid,
		TargetUserMids: SliceToSet(targets),
	}
	_, err := cl.conn.InviteIntoChatAsync(cl.client.ctx, req)
	return err
}
func (cl *TalkService) ReissueChatTicket(gid string) (string, error) {
	req := &model.ReissueChatTicketRequest{
		ReqSeq:   cl.client.RequestSequence,
		GroupMid: gid,
	}
	ticket, err := cl.conn.ReissueChatTicket(cl.client.ctx, req)
	if err != nil {
		return "", cl.client.afterError(err)
	}
	return ticket.TicketId, nil
}
func (cl *TalkService) ReissueChatTicketAsync(gid string) (<-chan *model.ReissueChatTicketResponse, <-chan error) {
	req := &model.ReissueChatTicketRequest{
		ReqSeq:   cl.client.RequestSequence,
		GroupMid: gid,
	}
	ticket, err := cl.conn.ReissueChatTicketAsync(cl.client.ctx, req)
	return ticket, err
}
func (cl *TalkService) RejectChatInvitation(gid string) error {
	req := &model.RejectChatInvitationRequest{
		ReqSeq:  cl.client.RequestSequence,
		ChatMid: gid,
	}
	err := cl.conn.RejectChatInvitation(cl.client.ctx, req)
	return cl.client.afterError(err)
}
func (cl *TalkService) UpdateChatName(gid, name string) error {
	chat := &model.Chat{ChatMid: gid, ChatName: name}
	req := &model.UpdateChatRequest{
		ReqSeq:           cl.client.RequestSequence,
		Chat:             chat,
		UpdatedAttribute: model.UpdatedChatAttribute_NAME,
	}
	_, err := cl.conn.UpdateChat(cl.client.ctx, req)
	return cl.client.afterError(err)
}
func (cl *TalkService) UpdateChatURL(chatID string, typeVar bool) error {
	if typeVar {
		return cl.closeChatUrlManual(chatID)
	}
	return cl.openChatUrlManual(chatID)
}
func (cl *TalkService) closeChatUrlManual(id string) error {
	request, _ := http.NewRequest("POST", PATH_NORMAL.ToURL(), bytes.NewBuffer([]byte("\x82!\x00\nupdateChat\x1c\x15\x00\x1c(!"+id+"l\x1c!\x00\x00\x00\x15\x08\x00\x00")))
	for key, value := range cl.client.ThriftFactory.header() {
		request.Header.Set(key, value)
	}
	_, err := cl.client.ThriftFactory.HttpClient().Do(request)
	return cl.client.afterError(err)
}

func (cl *TalkService) openChatUrlManual(id string) error {
	request, _ := http.NewRequest("POST", PATH_NORMAL.ToURL(), bytes.NewBuffer([]byte("\x82!\x00\nupdateChat\x1c\x15\x00\x1c(!"+id+"l\x1c\x00\x00\x00\x15\x08\x00\x00")))
	for key, value := range cl.client.ThriftFactory.header() {
		request.Header.Set(key, value)
	}
	_, err := cl.client.ThriftFactory.HttpClient().Do(request)
	return cl.client.afterError(err)
}

func (cl *TalkService) DeleteOtherFromChat(gid, mid string) error {
	req := &model.DeleteOtherFromChatRequest{
		ReqSeq:         cl.client.RequestSequence,
		ChatMid:        gid,
		TargetUserMids: map[string]bool{mid: true},
	}
	_, err := cl.conn.DeleteOtherFromChat(cl.client.ctx, req)
	return cl.client.afterError(err)
}
func (cl *TalkService) DeleteOtherFromChatAsync(gid, mid string) <-chan error {
	req := &model.DeleteOtherFromChatRequest{
		ReqSeq:         cl.client.RequestSequence,
		ChatMid:        gid,
		TargetUserMids: map[string]bool{mid: true},
	}
	_, err := cl.conn.DeleteOtherFromChatAsync(cl.client.ctx, req)
	return err
}
func (cl *TalkService) DeleteSelfFromChat(gid string) error {
	req := &model.DeleteSelfFromChatRequest{
		ReqSeq:  cl.client.RequestSequence,
		ChatMid: gid,
	}
	_, err := cl.conn.DeleteSelfFromChat(cl.client.ctx, req)
	return cl.client.afterError(err)
}
func (cl *TalkService) GetAllChatMids() (*model.GetAllChatMidsResponse, error) {
	req := &model.GetAllChatMidsRequest{
		WithMemberChats:  true,
		WithInvitedChats: true,
	}
	res, err := cl.conn.GetAllChatMids(cl.client.ctx, req, model.SyncReason_OPERATION)
	return res, cl.client.afterError(err)
}
func (cl *TalkService) CancelChatInvitation(gid, mid string) error {
	req := &model.CancelChatInvitationRequest{
		ReqSeq:         cl.client.RequestSequence,
		ChatMid:        gid,
		TargetUserMids: map[string]bool{mid: true},
	}
	_, err := cl.conn.CancelChatInvitation(cl.client.ctx, req)
	return cl.client.afterError(err)
}
func (cl *TalkService) FindChatByTicket(ticket string) (*model.Chat, error) {
	req := &model.FindChatByTicketRequest{
		TicketId: ticket,
	}
	chat, err := cl.conn.FindChatByTicket(cl.client.ctx, req)
	if err != nil {
		return nil, cl.client.afterError(err)
	}
	return chat.Chat, nil
}

/*
Contact functions
*/

func (cl *TalkService) FindAndAddContactByMid(mid string) error {
	_, err := cl.conn.FindAndAddContactsByMid(
		cl.client.ctx, cl.client.RequestSequence, mid,
		model.ContactType_MID, "",
	)
	return cl.client.afterError(err)
}
func (cl *TalkService) FindAndAddContactsByPhone(phones []string) (map[string]*model.Contact, error) {
	cons, err := cl.conn.FindAndAddContactsByPhone(
		cl.client.ctx, cl.client.RequestSequence, SliceToSet(phones), "",
	)
	return cons, cl.client.afterError(err)
}

func (cl *TalkService) GetContacts(mids []string) ([]*model.Contact, error) {
	contacts, err := cl.conn.GetContacts(cl.client.ctx, mids)
	return contacts, cl.client.afterError(err)
}
func (cl *TalkService) GetContact(mid string) (*model.Contact, error) {
	contact, err := cl.conn.GetContact(cl.client.ctx, mid)
	return contact, cl.client.afterError(err)
}
func (cl *TalkService) BlockContact(mid string) error {
	return cl.client.afterError(cl.conn.BlockContact(cl.client.ctx, cl.client.RequestSequence, mid))
}

func (cl *TalkService) UpdateContactSetting(mid, attr string, val model.UpdateContactSettingFlag) error {
	err := cl.conn.UpdateContactSetting(cl.client.ctx, cl.client.RequestSequence, mid, val, attr)
	return cl.client.afterError(err)
}
func (cl *TalkService) GetAllContactIds() ([]string, error) {
	res, err := cl.conn.GetAllContactIds(cl.client.ctx, model.SyncReason_OPERATION)
	return res, cl.client.afterError(err)
}
func (cl *TalkService) GetRecommendationIds() {
	_, _ = cl.conn.GetRecommendationIds(cl.client.ctx, model.SyncReason_INITIALIZATION)
}
func (cl *TalkService) GetBlockedContactIds() {
	_, _ = cl.conn.GetBlockedContactIds(cl.client.ctx, model.SyncReason_INITIALIZATION)
}
func (cl *TalkService) GetBlockedRecommendationIds() {
	_, _ = cl.conn.GetBlockedRecommendationIds(cl.client.ctx, model.SyncReason_INITIALIZATION)
}
func (cl *TalkService) FindContactByTicket(ticket string) (*model.Contact, error) {
	contact, err := cl.conn.FindContactByUserTicket(cl.client.ctx, ticket)
	return contact, cl.client.afterError(err)
}
func (cl *TalkService) FindContactByUserId(id string) (*model.Contact, error) {
	contact, err := cl.conn.FindContactByUserid(cl.client.ctx, id)
	return contact, cl.client.afterError(err)
}
func (cl *TalkService) AddContactUsingTicket(ticket string) error {
	contact, err := cl.FindContactByTicket(ticket)
	if err != nil {
		return cl.client.afterError(err)
	}
	_, err = cl.conn.FindAndAddContactsByMid(
		cl.client.ctx, cl.client.RequestSequence, contact.Mid,
		model.ContactType_MID, "{\"screen\":\"urlScheme:internal\",\"spec\":\"native\",\"ticketId\":\""+ticket+"\"}",
	)
	return cl.client.afterError(err)
}
func (cl *TalkService) AddContactByUserId(id string) error {
	_, err := cl.FindContactByUserId(id)
	if err != nil {
		return cl.client.afterError(err)
	}
	_, err = cl.conn.FindAndAddContactsByUserid(cl.client.ctx, cl.client.RequestSequence, id, "{\"screen\":\"friendAdd:idSearch\",\"spec\":\"native\"}")
	return cl.client.afterError(err)
}
func (cl *TalkService) AddContactGroupMember(mid string) error {
	_, err := cl.conn.FindAndAddContactsByMid(
		cl.client.ctx, cl.client.RequestSequence, mid,
		model.ContactType_MID, "{\"screen\":\"groupMemberList\",\"spec\":\"native\"}",
	)
	return cl.client.afterError(err)
}

/*
Setting functions
*/

func (cl *TalkService) GetSettings(reason model.SyncReason) (*model.Settings, error) {
	settings, err := cl.conn.GetSettings(cl.client.ctx, reason)
	if err != nil {
		cl.client.Settings = settings
	}
	return settings, cl.client.afterError(err)
}
func (cl *TalkService) UpdateSettingsAttributes2(attributesToUpdate map[model.PendingAgreement]bool, settings *model.Settings) error {
	_, err := cl.conn.UpdateSettingsAttributes2(cl.client.ctx, cl.client.RequestSequence, attributesToUpdate, settings)
	return cl.client.afterError(err)
}

/*
Other functions
*/

func (cl *TalkService) Noop() (err error) {
	err = cl.conn.Noop(cl.client.ctx)
	return cl.client.afterError(err)
}
func (cl *TalkService) GetPendingAgreements() ([]model.PendingAgreement, error) {
	agreements, err := cl.conn.GetPendingAgreements(cl.client.ctx)
	if agreements == nil {
		return nil, cl.client.afterError(err)
	}
	return agreements.PendingAgreements, cl.client.afterError(err)
}
func (cl *TalkService) GetConfigurations(reason model.SyncReason) error {
	_, err := cl.conn.GetConfigurations(cl.client.ctx, 0, "JP", cl.client.ClientInfo.PhoneNumber.CountryCode, "JP", "44010", reason)
	return cl.client.afterError(err)
}

func (cl *TalkService) NotifyRegistrationComplete() error {
	err := cl.conn.NotifyRegistrationComplete(cl.client.ctx, cl.client.ClientInfo.Device.Udid, cl.client.GetLineApplicationHeader())
	return cl.client.afterError(err)
}

func (cl *TalkService) FollowUser(mid string) error {
	return cl.client.afterError(cl.conn.Follow(cl.client.ctx, &model.FollowRequest{FollowMid: &model.FollowMid{Mid: mid}}))
}

func (cl *TalkService) UnFollowUser(mid string) error {
	return cl.client.afterError(cl.conn.Unfollow(cl.client.ctx, &model.UnfollowRequest{FollowMid: &model.FollowMid{Mid: mid}}))
}

func (cl *TalkService) CreateRoom(mids []string) (*model.Room, error) {
	room, err := cl.conn.CreateRoomV2(cl.client.ctx, cl.client.RequestSequence, mids)
	return room, cl.client.afterError(err)
}

func (cl *TalkService) SyncContacts(cons []*model.ContactModification) (map[string]*model.ContactRegistration, error) {
	res, err := cl.conn.SyncContacts(cl.client.ctx, cl.client.RequestSequence, cons)
	return res, cl.client.afterError(err)
}

func (cl *TalkService) GetContactsV2(mids []string) (*model.GetContactsV2Response, error) {
	tmp := true
	req := &model.GetContactsV2Request{
		TargetUserMids: mids,
		NeededContactCalendarEvents: map[model.ContactCalendarEventType]bool{
			model.ContactCalendarEventType_BIRTHDAY: true,
		},
		WithUserStatus: &tmp,
	}
	res, err := cl.conn.GetContactsV2(cl.client.ctx, req, model.SyncReason_INITIALIZATION)
	return res, cl.client.afterError(err)
}
