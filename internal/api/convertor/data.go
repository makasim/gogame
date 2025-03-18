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

func FindGame(e flowstate.Engine, gID string, gRev int32) (*v1.Game, *flowstate.StateCtx, *flowstate.Data, error) {
	d := &flowstate.Data{}
	stateCtx := &flowstate.StateCtx{}

	if err := e.Do(
		flowstate.GetByID(stateCtx, flowstate.StateID(gID), int64(gRev)),
		flowstate.DereferenceData(stateCtx, d, `game`),
		flowstate.GetData(d),
	); err != nil {
		return nil, nil, nil, err
	}

	g, err := DataToGame(d)
	if err != nil {
		return nil, nil, nil, err
	}

	g.Rev = int32(stateCtx.Current.Rev)

	return g, stateCtx, d, nil
}
