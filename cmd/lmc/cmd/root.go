package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ollybritton/go-lmc"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "lmc",
	Short: "Simulate a little man computer",
	Long:  `Simulate a little man computer`,
	Args:  cobra.ExactArgs(1),

	Run: func(cmd *cobra.Command, args []string) {
		filename := args[0]

		opcodeSize, err := cmd.Flags().GetInt("opcode-size")
		checkFlagErr(err)
		operandSize, err := cmd.Flags().GetInt("operand-size")
		checkFlagErr(err)
		shouldStep, err := cmd.Flags().GetBool("step")
		checkFlagErr(err)
		shouldLog, err := cmd.Flags().GetBool("log")
		checkFlagErr(err)

		if shouldLog {
			logrus.SetLevel(logrus.DebugLevel)
		}

		bytes, err := ioutil.ReadFile(filename)
		if err != nil {
			logrus.Errorf("Error reading file: %s", err)
		}

		computer, err := lmc.NewComputerFromCode(string(bytes), opcodeSize, operandSize)
		if err != nil {
			logrus.Fatal(err)
		}

		go func() {
			err = computer.Run()
			if err != nil {
				logrus.Errorf("Error running computer: %s", err)
			}
		}()

		for {
			msg := <-computer.Messages
			switch msg.Status {
			case lmc.Done:
				logrus.Info("DONE")
				return
			case lmc.NeedInput:
				logrus.Infoln("NEED INPUT")
				var i int
				fmt.Print("Input (int): ")
				fmt.Scan(&i)
				computer.Inbox <- i
			case lmc.NeedStep:
				if shouldStep {
					logrus.Infoln("NEED STEP")
					fmt.Scanln()
				}
				computer.Step <- struct{}{}
			case lmc.Log:
				logrus.Debugln(msg.Val)
			case lmc.Output:
				logrus.Infoln("OUTPUT", msg.Val)
			}
		}
	},
}

func init() {
	rootCmd.Flags().IntP("opcode-size", "c", 1, "size of opcode, in digits")
	rootCmd.Flags().IntP("operand-size", "r", 2, "size of operand, in digits")
	rootCmd.Flags().BoolP("step", "s", false, "whether to step through the input")

	rootCmd.Flags().BoolP("log", "l", false, "whether to log each stage of the computation")

}

func checkFlagErr(err error) {
	if err != nil {
		logrus.Fatalf("Error getting flag: %s", err)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
