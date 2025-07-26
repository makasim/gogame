package convertor

import (
	"github.com/makasim/flowstate"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

func UndoToData(u *v1.Undo, d *flowstate.Data) error {
	b, err := protojson.Marshal(u)
	if err != nil {
		return err
	}

	//d.ID = flowstate.DataID(fmt.Sprintf("%s-%d", u.GameId, u.GameRev))
	d.Blob = b
	return nil
}

func DataToUndo(d *flowstate.Data) (*v1.Undo, error) {
	u := &v1.Undo{}
	if err := protojson.Unmarshal(d.Blob, u); err != nil {
		return nil, err
	}

	return u, nil
}
