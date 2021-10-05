package line

import (
	"github.com/line-api/model/go/model"
	"strconv"
	"strings"
)

type PollData struct {
	LastRev       int64
	Count         int32
	GlobalRev     int64
	IndividualRev int64
}

type PollService struct {
	client *Client

	conn     *model.FTalkServiceClient
	connTMCP *model.FTalkServiceClient
	PollData *PollData
}

func (cl *Client) newPollService() *PollService {
	return &PollService{
		client:   cl,
		conn:     cl.thriftFactory.newPollServiceClient(),
		connTMCP: cl.thriftFactory.newPollTMCPServiceClient(),
		PollData: &PollData{
			Count: 50,
		},
	}
}

func (s *PollService) FetchLineOperations() ([]*model.Operation, error) {
	return s.fetchLineOperationsInternal(s.FetchOps)
}

func (s *PollService) FetchLineOperationsTMCP() ([]*model.Operation, error) {
	return s.fetchLineOperationsInternal(s.FetchOpsTMCP)
}

func (s *PollService) fetchLineOperationsInternal(getter func() ([]*model.Operation, error)) ([]*model.Operation, error) {
	ops, err := getter()
	if err != nil {
		return nil, err
	}
	var operations []*model.Operation
	for _, op := range ops {
		if op.OpType == model.OpType_END_OF_OPERATION {
			if op.Param2 != "" {
				s.PollData.GlobalRev = s.getGlobalRev(op)
			}
			if op.Param1 != "" {
				s.PollData.IndividualRev = s.getIndividualRev(op)
			}
		} else {
			operations = append(operations, op)
		}
		s.setRevision(op.Revision)
	}
	return operations, nil
}

func (s *PollService) FetchOps() ([]*model.Operation, error) {
	ops, err := s.conn.FetchOps(
		s.client.ctx,
		s.PollData.LastRev, s.PollData.Count,
		s.PollData.GlobalRev, s.PollData.IndividualRev,
	)
	return ops, s.client.afterError(err)
}
func (s *PollService) FetchOpsTMCP() ([]*model.Operation, error) {
	ops, err := s.connTMCP.FetchOps(
		s.client.ctx,
		s.PollData.LastRev, s.PollData.Count,
		s.PollData.GlobalRev, s.PollData.IndividualRev,
	)
	return ops, s.client.afterError(err)
}

func (s *PollService) FetchOperations() ([]*model.Operation, error) {
	ops, err := s.conn.FetchOperations(s.client.ctx,
		s.PollData.LastRev, s.PollData.Count,
	)
	return ops, s.client.afterError(err)
}

func (s *PollService) getIndividualRev(op *model.Operation) int64 {
	if op.Param1 != "" {
		sps := strings.Split(op.Param1, "")
		if len(sps) != 0 {
			res, _ := strconv.ParseInt(sps[0], 10, 64)
			return res
		}
	}
	return 0
}

func (s *PollService) getGlobalRev(op *model.Operation) int64 {
	if op.Param2 != "" {
		sps := strings.Split(op.Param2, "")
		if len(sps) != 0 {
			res, _ := strconv.ParseInt(sps[0], 10, 64)
			return res
		}
	}
	return 0
}

func (s *PollService) setRevision(rev int64) {
	if s.PollData.LastRev < rev {
		s.PollData.LastRev = rev
	}
}
