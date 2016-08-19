/*
Copyright (c) 2014 Ashley Jeffs

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package output

import (
	"errors"

	"github.com/jeffail/benthos/lib/broker"
	"github.com/jeffail/benthos/lib/types"
	"github.com/jeffail/util/log"
	"github.com/jeffail/util/metrics"
)

//--------------------------------------------------------------------------------------------------

var (
	// ErrFanOutNoOutputs - Returned when creating a FanOut type with zero outputs.
	ErrFanOutNoOutputs = errors.New("attempting to create fan_out output type with no outputs")
)

//--------------------------------------------------------------------------------------------------

func init() {
	constructors["fan_out"] = typeSpec{
		constructor: NewFanOut,
		description: `
The 'fan_out' output type allows you to group multiple outputs together. With
the fan out model all outputs will receive every message that passes through
benthos. This process is blocking, meaning if any output applies backpressure
then it will block all outputs from receiving messages.`,
	}
}

//--------------------------------------------------------------------------------------------------

// FanOutConfig - Configuration for the FanOut output type.
type FanOutConfig struct {
	Outputs []interface{} `json:"outputs" yaml:"outputs"`
}

// NewFanOutConfig - Creates a new FanOutConfig with default values.
func NewFanOutConfig() FanOutConfig {
	return FanOutConfig{
		Outputs: []interface{}{},
	}
}

//--------------------------------------------------------------------------------------------------

/*
NewFanOut - Create a new FanOut output type. Messages will be sent out to ALL outputs, outputs which
block will apply backpressure upstream, meaning other outputs will also stop receiving messages.
*/
func NewFanOut(conf Config, log log.Modular, stats metrics.Aggregator) (Type, error) {
	if len(conf.FanOut.Outputs) == 0 {
		return nil, ErrFanOutNoOutputs
	}

	outputConfs, err := parseOutputConfsWithDefaults(conf.FanOut.Outputs)
	if err != nil {
		return nil, err
	}

	outputs := make([]types.Consumer, len(outputConfs))

	for i, oConf := range outputConfs {
		outputs[i], err = Construct(oConf, log, stats)
		if err != nil {
			return nil, err
		}
	}

	return broker.NewFanOut(outputs, log, stats)
}

//--------------------------------------------------------------------------------------------------