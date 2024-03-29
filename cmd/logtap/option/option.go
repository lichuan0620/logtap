package option

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	flag "github.com/spf13/pflag"

	"github.com/lichuan0620/logtap/cmd/logtap/version"
	"github.com/lichuan0620/logtap/pkg/fieldpath"
	model "github.com/lichuan0620/logtap/pkg/model/v1alpha1"
	"github.com/lichuan0620/logtap/pkg/signal"
)

const (
	defaultTimestamp  = time.RFC3339Nano
	defaultMinSize    = 128
	defaultInterval   = 0.5
	defaultWebAddress = ":8080"

	noDefault = ""
)

var (
	// Spec is the LogTaskSpec created according the the command line options.
	Spec = new(model.LogTaskSpec)

	// WebAddress is the address to listen on for most HTTP requests.
	WebAddress string

	// Name is used to differentiate different deployments.
	Name string

	// StopCh closes when the program should be cleaned and terminated.
	StopCh chan struct{}
)

var (
	commandLine = flag.NewFlagSet(version.Name, flag.ExitOnError)

	template = commandLine.StringP(
		"template", "t", getEnv("LOGTAP_TEMPLATE", noDefault),
		"The name of the predefined template to run; override other options if specified",
	)

	duration = commandLine.String(
		"duration", getEnv("LOGTAP_DURATION", noDefault),
		"The duration, such as 1h or 30s, for which LogTap should run",
	)

	outputKindHelp = []string{
		fmt.Sprintf(
			"  %s\tThe log messages will be written to STDOUT",
			model.OutputKindStdOut,
		),
		fmt.Sprintf(
			"  %s\tThe log messages will be written to STDERR",
			model.OutputKindStdErr,
		),
		fmt.Sprintf(
			"  %s\t\tThe log messages will be written to the specified file",
			model.OutputKindFile,
		),
	}

	contentTypeHelp = []string{
		fmt.Sprintf(
			"  %s\tThe log messages will be randomly generated with a minimal size",
			model.ContentTypeRandom,
		),
		fmt.Sprintf(
			"  %s\tThe log messages will be explicitly defined",
			model.ContentTypeExplicit,
		),
	}

	presetHelp = []string{
		fmt.Sprintf(
			"  %s\tProduces a load of 256 B/log, 10 logs/s, and 2.5 KiB/s",
			model.TaskPresetStandard,
		),
		fmt.Sprintf(
			"  %s\t\tProduces a load of 20 MiB/log, 0.5 log/s, and 10 Mib/s",
			model.TaskPresetLong,
		),
		fmt.Sprintf(
			"  %s\tProduces a load of 256 B/log, 50000 log/s, and 12 Mib/s",
			model.TaskPresetFrequent,
		),
		fmt.Sprintf(
			"  %s\t\tProduces a load of 1 MiB/log, 40 log/s, and 40 Mib/s",
			model.TaskPresetRoast,
		),
	}

	extraHelp = fmt.Sprintf(`
Output Kinds:
%s

Content Types:
%s

Task Presets:
%s`,
		strings.Join(outputKindHelp, "\n"),
		strings.Join(contentTypeHelp, "\n"),
		strings.Join(presetHelp, "\n"),
	)

	usage = fmt.Sprintf(`Logtap is a benchmark tool that generates log messages in a controlled way.

Find more information at https://github.com/lichuan0620/logtap

Options:`)
)

func init() {
	flag.ErrHelp = fmt.Errorf("")
	commandLine.Usage = printHelp
	parse()
	failOnError(model.ValidateLogTaskSpec(fieldpath.NewFieldPath(), Spec))
}

func parse() {

	commandLine.StringVarP(
		&Name, "name", "n", getEnv("LOGTAP_NAME", getEnv("HOST", version.Name)),
		"The name given to LogTap to differentiate different deployments",
	)

	commandLine.StringVar(
		&WebAddress, "web.address", getEnv("LOGTAP_WEB_ADDRESS", defaultWebAddress),
		"The address to listen on for most HTTP requests",
	)

	commandLine.StringVar(&Spec.OutputKind,
		"output.kind", getEnv("LOGTAP_OUTPUT_KIND", model.OutputKindStdErr),
		"The channel to which the log messages should be sent",
	)

	commandLine.StringVar(&Spec.Filepath,
		"output.filePath", getEnv("LOGTAP_OUTPUT_FILE_PATH", noDefault),
		"Path to the log file to which the log messages would be appended",
	)

	commandLine.StringVar(&Spec.TimestampFormat,
		"timestamp.format", getEnv("LOGTAP_TIMESTAMP_FORMAT", defaultTimestamp),
		"Format of the log timestamp",
	)

	timestampOff := commandLine.Bool(
		"timestamp.off", getBoolEnv("LOGTAP_TIMESTAMP_OFF", false),
		"Disable log timestamp",
	)

	commandLine.StringVar(&Spec.ContentType,
		"content.type", getEnv("LOGTAP_CONTENT_TYPE", model.ContentTypeRandom),
		"The type of content that the log messages would have",
	)

	commandLine.StringVar(&Spec.Message,
		"content.message", getEnv("LOGTAP_CONTENT_MESSAGE", noDefault),
		"The log message to be be printed",
	)

	commandLine.IntVarP(&Spec.MinSize,
		"content.minSize", "s", getIntEnv("LOGTAP_CONTENT_MIN_SIZE", defaultMinSize),
		"The minimal size of a randomized log message in bytes",
	)

	commandLine.Float64VarP(&Spec.Interval,
		"interval", "i", getFloat64Env("LOGTAP_INTERVAL", defaultInterval),
		"The amount of time, in seconds, to wait in-between log messages",
	)

	showVersion := commandLine.BoolP(
		"version", "v", false,
		"Print the version information and quit",
	)

	showHelp := commandLine.BoolP(
		"help", "h", false,
		"Print the help information and quit",
	)

	commandLine.Parse(os.Args[1:])

	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	if *timestampOff {
		Spec.TimestampFormat = ""
	}

	if len(*template) > 0 {
		var err error
		Spec, err = model.GetLogTaskSpecPreset(*template)
		failOnError(err)
	}

	if len(*duration) > 0 {
		d, err := time.ParseDuration(*duration)
		failOnError(err)
		StopCh = make(chan struct{})
		go terminateAfter(d)
	} else {
		StopCh = signal.SetupStopSignalHandler()
	}
}

func printVersion() {
	fmt.Fprintln(os.Stderr, fmt.Sprintf("%s version %s", version.Name, version.Version))
}

func printHelp() {
	fmt.Fprintln(os.Stderr, usage)
	commandLine.PrintDefaults()
	fmt.Fprintln(os.Stderr, extraHelp)
}

func failOnError(err error) {
	if err != nil {
		printHelp()
		fmt.Fprintf(os.Stderr, "\n%s\n", err.Error())
		os.Exit(2)
	}
}

func terminateAfter(duration time.Duration) {
	select {
	case <-time.After(duration):
	case <-signal.SetupStopSignalHandler():
	}
	close(StopCh)
}

func getEnv(name, def string) string {
	if env := os.Getenv(name); env != "" {
		return env
	}
	return def
}

func getBoolEnv(name string, def bool) bool {
	if env := os.Getenv(name); env != "" {
		return strings.ToLower(env) == "true"
	}
	return def
}

func getIntEnv(name string, def int) int {
	if env := os.Getenv(name); env != "" {
		if ret, err := strconv.Atoi(env); err == nil {
			return ret
		}
	}
	return def
}

func getFloat64Env(name string, def float64) float64 {
	if env := os.Getenv(name); env != "" {
		if ret, err := strconv.ParseFloat(env, 64); err == nil {
			return ret
		}
	}
	return def
}
