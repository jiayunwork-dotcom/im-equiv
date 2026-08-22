// Command im-equiv is a command-line calculator for the T-equivalent circuit
// of a three-phase induction machine. Given the per-phase circuit parameters
// and the slip, it computes the input current, the electromagnetic torque,
// the mechanical power and the slip of maximum torque, together with the
// supporting power split and speeds.
//
// Usage:
//
//	go run . torque example/4pole-50hz.json
//
// The torque subcommand reads the JSON input file, validates it and prints
// the computed quantities. Any invalid input is reported on standard error
// and the process exits with a non-zero status: exit code 2 for a usage
// error (wrong arguments or an unknown subcommand) and exit code 1 for a run
// error (unreadable file, malformed JSON, or a physically invalid input).
//
// The input file follows the schema documented in package config. The
// shipped example is a four-pole, 50 Hz, 380 V machine at 3% slip; a second
// example, example/4pole-50hz-generating.json, runs the same machine as a
// generator at -2% slip.
//
// The command prints an aligned table of quantities. The essential readings
// are the input current I1, the electromagnetic torque T, the mechanical
// power Pem and the maximum-torque slip sm; the surrounding rows (speeds,
// air-gap power, copper losses, efficiency) let the user verify the energy
// balance Pag = Pcu2 + Pem by eye.
package main

import (
	"fmt"
	"io"
	"os"

	"im-equiv/internal/calc"
	"im-equiv/internal/config"
)

// exit codes follow the convention that usage errors are 2 and run errors 1.
const (
	exitUsage = 2
	exitRun   = 1
)

func main() {
	if len(os.Args) < 2 {
		usage(os.Stderr)
		os.Exit(exitUsage)
	}

	switch os.Args[1] {
	case "torque":
		if err := runTorque(os.Args[2:], os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, "im-equiv: error:", err)
			os.Exit(exitRun)
		}
	case "help", "-h", "--help":
		usage(os.Stdout)
	default:
		fmt.Fprintf(os.Stderr, "im-equiv: unknown subcommand %q\n", os.Args[1])
		usage(os.Stderr)
		os.Exit(exitUsage)
	}
}

// runTorque loads the input file, computes and prints the results. It is
// separated from main so the error paths can be tested without terminating
// the test process. The output writer is injected for the same reason.
func runTorque(args []string, out io.Writer) error {
	if len(args) != 1 {
		return fmt.Errorf("torque requires exactly one JSON input file, got %d arguments", len(args))
	}
	in, err := config.LoadFile(args[0])
	if err != nil {
		return err
	}
	res, err := calc.Compute(in)
	if err != nil {
		return err
	}
	_, err = io.WriteString(out, calc.FormatResults(res))
	return err
}

func usage(w io.Writer) {
	fmt.Fprintln(w, "im-equiv: three-phase induction machine T-equivalent circuit calculator")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "usage:")
	fmt.Fprintln(w, "  im-equiv torque <input.json>")
	fmt.Fprintln(w, "      solve one operating point and print currents, torque,")
	fmt.Fprintln(w, "      mechanical power, max-torque slip and the power split")
	fmt.Fprintln(w, "  im-equiv help")
	fmt.Fprintln(w, "      print this help text")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "exit codes:")
	fmt.Fprintln(w, "  0  success")
	fmt.Fprintln(w, "  1  run error: unreadable file, malformed JSON or")
	fmt.Fprintln(w, "     a physically invalid input")
	fmt.Fprintln(w, "  2  usage error: unknown subcommand or wrong arguments")
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, "example:")
	fmt.Fprintln(w, "  im-equiv torque example/4pole-50hz.json")
}
