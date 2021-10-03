package line

import (
	"encoding/json"
	"fmt"
	"golang.org/x/xerrors"
	"io/ioutil"
	"os"
)

func CreateDirIfNotExist(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		err := os.Mkdir(dirName, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func isPathExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func ReadJsonToStruct(fName string, struct_ interface{}) (interface{}, error) {
	file, err := os.Open(fName)
	if err != nil {
		return nil, err
	}
	parser := json.NewDecoder(file)
	err = parser.Decode(struct_)
	return struct_, err
}

func WriteStructToJson(fName string, struct_ interface{}) error {
	file, err := json.MarshalIndent(struct_, "", "    ")
	err = ioutil.WriteFile(fName, file, 0644)
	return err
}

func (cl *Client) SaveKeeper() error {
	err := CreateDirIfNotExist(cl.ClientSetting.KeeperDir)
	if err != nil {
		return err
	}
	path := fmt.Sprintf(cl.ClientSetting.KeeperDir+"/%v.keeper", cl.Profile.Mid)
	return WriteStructToJson(path, cl)
}

func (cl *Client) LoadKeeper() error {
	path := fmt.Sprintf(cl.ClientSetting.KeeperDir+"/%v.keeper", cl.Profile.Mid)
	if isPathExist(path) {
		_, err := ReadJsonToStruct(path, cl)
		return err
	}
	return xerrors.New("no keeper file found")
}
