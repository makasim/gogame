package convertor

import (
	"github.com/makasim/flowstate"
	v1 "github.com/makasim/gogame/protogen/gogame/v1"
	"google.golang.org/protobuf/encoding/protojson"
)

func GameToData(g *v1.Game, d *flowstate.Data) error {
	b, err := protojson.Marshal(g)
	if err != nil {
		return err
	}

	d.ID = flowstate.DataID(g.Id)
	d.B = b
	return nil
}

func DataToGame(d *flowstate.Data) (*v1.Game, error) {
	g := &v1.Game{}
	if err := protojson.Unmarshal(d.B, g); err != nil {
		return nil, err
	}

	return g, nil
}
