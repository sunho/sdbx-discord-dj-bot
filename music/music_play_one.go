package music

import (
	"bufio"
	"encoding/binary"
	"log"
	"os/exec"

	"github.com/bwmarrin/discordgo"
	"github.com/sunho/sdbx-discord-dj-bot/music/provider"
)

type Song struct {
	provider.Song
	Requestor *discordgo.Member
}

func playOne(conn *discordgo.VoiceConnection, bufferSize int, skipC chan struct{}, url string) {
	ytdl := exec.Command("youtube-dl", "-v", "-f", "bestaudio", "-o", "-", url)
	ytdlout, err := ytdl.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}

	ffmpeg := exec.Command("ffmpeg", "-i", "pipe:0", "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
	ffmpegout, err := ffmpeg.StdoutPipe()
	ffmpeg.Stdin = ytdlout
	if err != nil {
		log.Println(err)
		return
	}
	ffmpegbuf := bufio.NewReaderSize(ffmpegout, bufferSize)

	dca := exec.Command("dca")
	dca.Stdin = ffmpegbuf
	dcaout, err := dca.StdoutPipe()
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		go dca.Wait()
	}()
	dcabuf := bufio.NewReaderSize(dcaout, bufferSize)

	log.Println(bufferSize)

	err = ytdl.Start()
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		go ytdl.Wait()
	}()

	err = ffmpeg.Start()
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		go ffmpeg.Wait()
	}()

	err = dca.Start()
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		go dca.Wait()
	}()

	if dcabuf == nil {
		return
	}

	var opuslen int16
	conn.Speaking(true)
	defer conn.Speaking(false)
	for {
		select {
		case <-skipC:
			return
		default:
			err = binary.Read(dcabuf, binary.LittleEndian, &opuslen)
			if err != nil {
				log.Println(err)
				return
			}

			opus := make([]byte, opuslen)
			err = binary.Read(dcabuf, binary.LittleEndian, &opus)
			if err != nil {
				log.Println(err)
				return
			}

			conn.OpusSend <- opus
		}
	}
}
