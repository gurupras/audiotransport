package audiotransport

import (
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/alecthomas/kingpin"
)

var (
	name       *string
	addr       *string
	proto      *string
	device     *string
	api        *string
	verbose    *bool
	samplerate *uint32
	channels   *uint32
)

type ArgsConfig struct {
	Api        ApiType
	Proto      string
	Name       string
	Device     string
	Samplerate uint32
	Channels   uint32
	Addrs      []string
}

func setupParser(app *kingpin.Application) {
	name = app.Flag("name", "program name. This is used as filename in FILE method").Short('n').Required().String()
	addr = app.Arg("receiver-address", "Address of receiver").Required().String()
	proto = app.Flag("protocol", "tcp/udp").Short('P').Default("udp").String()
	device = app.Flag("device", "Device from which to capture and transmit").Short('d').String()
	api = app.Flag("method", "Which mechanism to use.. ALSA/PULSE").Short('m').Default("PULSE").String()
	samplerate = app.Flag("samplerate", "The samplerate to use").Short('s').Default("44100").Uint32()
	channels = app.Flag("channels", "Number of channels to use").Short('c').Default("2").Uint32()
	verbose = app.Flag("verbose", "Enable verbose logging").Short('v').Default("false").Bool()

}

func ParseArgs(app *kingpin.Application, args []string) *ArgsConfig {
	setupParser(app)
	kingpin.MustParse(app.Parse(args[1:]))

	if *verbose {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Enabling verbose logging")
	}

	var apiType ApiType
	switch *api {
	case "PULSE":
		apiType = PULSE_API
	case "ALSA":
		apiType = ALSA_API
	case "FILE":
		apiType = FILE_API
	default:
		log.Fatalf("Invalid API: %v", *api)
	}
	log.Infof("Using API: %v", apiType.ApiString())

	addrs := strings.Split(*addr, ",")

	config := &ArgsConfig{}
	config.Api = apiType
	config.Proto = *proto
	config.Name = *name
	config.Device = *device
	config.Samplerate = *samplerate
	config.Channels = *channels
	config.Addrs = addrs

	return config
}
