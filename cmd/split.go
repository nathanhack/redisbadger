package cmd

import (
	"bufio"
	"io"
	"os"

	"github.com/cheggaaa/pb"
	"github.com/nathanhack/aof"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// splitCmd represents the split command
var splitCmd = &cobra.Command{
	Use:   "split APPENDONLY.AOF",
	Short: "Takes an AOF file and splits it into two files: good and bad.",
	Long: `Takes an AOF file and outputs two files 
(APPENDONLY.AOF.good and APPENDONLY.AOF.bad), one containing valid 
commands and the other the bytes from that did not contain a command. 
A special note, commands are simple and may contain bad data as 
validation on the Key and Value must be done by the user some other way.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputFile, err := os.Open(args[0])
		if err != nil {
			logrus.Errorf("expected to find an AOF file: %v", err)
			return err
		}
		defer inputFile.Close()

		fixedPathname := args[0] + ".fixed"
		goodFile, err := os.Create(fixedPathname)
		if err != nil {
			logrus.Errorf("output file %v must NOT already exists: %v", err)
			return err
		}

		badPathname := args[0] + ".bad"
		badFile, err := os.Create(badPathname)
		if err != nil {
			logrus.Errorf("output file %v must NOT already exists: %v", err)
			return err
		}

		inputReader := bufio.NewReader(inputFile)
		goodWriter := bufio.NewWriter(goodFile)
		badWriter := bufio.NewWriter(badFile)

		splitter(inputReader, goodWriter, badWriter)

		err = goodWriter.Flush()
		if err != nil {
			logrus.Warn(err)
		}
		err = badWriter.Flush()
		if err != nil {
			logrus.Warn(err)
		}

		err = badFile.Sync()
		if err != nil {
			logrus.Warn(err)
		}
		err = goodFile.Sync()
		if err != nil {
			logrus.Warn(err)
		}

		err = badFile.Close()
		if err != nil {
			logrus.Warn(err)
		}
		err = goodFile.Close()
		if err != nil {
			logrus.Warn(err)
		}

		return nil
	},
}

func init() {
	aofCmd.AddCommand(splitCmd)
}

func splitter(inputReader *bufio.Reader, goodWriter, badWriter *bufio.Writer) {
	bar := pb.StartNew(0)
	run := true
	for op, buf, err := aof.ReadCommand(inputReader); run; op, buf, err = aof.ReadCommand(inputReader) {
		bar.Add(1)
		run = err != io.EOF
		if err == nil || (err == io.EOF && op != nil) {
			err = aof.WriteCommand(op, goodWriter)
			if err != nil {
				logrus.Error(err)
				return
			}
			continue
		}

		// well we had some sort of err that wasn't good
		// so we now write the current bytes out to
		// the "bad" writer
		n, err := badWriter.Write(buf)
		if err != nil {
			logrus.Error(err)
			return
		}
		if n != len(buf) {
			logrus.Errorf("write (%v) failed to match expected %v", len(buf), n)
			return
		}
	}

	bar.Finish()
}
