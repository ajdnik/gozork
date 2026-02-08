package game

import "testing"

func TestHarnessMetaMovementScript(t *testing.T) {
	runScript(t, []Step{
		{Command: "verbose", Contains: []string{"Maximum verbosity"}},
		{Command: "brief", Contains: []string{"Brief descriptions"}},
		{Command: "superbrief", Contains: []string{"Superbrief descriptions"}},
		{Command: "version", Contains: []string{"ZORK I"}},
		{Command: "verify", Contains: []string{"Verifying disk", "disk is correct"}},
		{Command: "bug", Contains: []string{"flawless program"}},
		{Command: "zork", Contains: []string{"At your service"}},
		{Command: "wish", Contains: []string{"wish will come true"}},
		{Command: "wait", Contains: []string{"Time passes"}},
		{Command: "inventory", Contains: []string{"empty-handed"}},
		{Command: "walk to mailbox", Contains: []string{"it's here"}},
		{Command: "north", Contains: []string{"North of House"}},
	})
}
