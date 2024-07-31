package convertor

import (
	"github.com/makasim/flowstate"
	v2 "github.com/makasim/gogame/protogen/gogame/v2"
	"google.golang.org/protobuf/encoding/protojson"
)

func GameToData(g *v2.Game, d *flowstate.Data) error {
	b, err := protojson.Marshal(g)
	if err != nil {
		return err
	}

	d.ID = flowstate.DataID(g.Id)
	d.B = b
	return nil
}

func DataToGame(d *flowstate.Data) (*v2.Game, error) {
	g := &v2.Game{}
	if err := protojson.Unmarshal(d.B, g); err != nil {
		return nil, err
	}

	return g, nil
}

func FindGame(e *flowstate.Engine, gID string, gRev int32) (*v2.Game, *flowstate.StateCtx, *flowstate.Data, error) {
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
