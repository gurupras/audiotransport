#include <stdio.h>
#include <unistd.h>
#include <alsa/asoundlib.h>

static snd_pcm_t *PlaybackHandle;

int init_playback(const char *device, int samplerate, int channels) {
	int err;

	printf("Init parameters: %s %d %d\n", device, samplerate, channels);

	if((err = snd_pcm_open(&PlaybackHandle, device, SND_PCM_STREAM_PLAYBACK, 0)) < 0) {
		printf("Can't open audio %s: %s\n", device, snd_strerror(err));
		return -1;
	}

	if ((err = snd_pcm_set_params(PlaybackHandle, SND_PCM_FORMAT_S16, SND_PCM_ACCESS_RW_INTERLEAVED, channels, samplerate, 1, 500000)) < 0) {
		printf("Can't set sound parameters: %s\n", snd_strerror(err));
		return -1;
	}
	return 0;
}

int alsa_writei(const void *bytes, int len) {
	snd_pcm_uframes_t frames = snd_pcm_writei(PlaybackHandle, bytes, len);
	if (frames < 0)
		frames = snd_pcm_recover(PlaybackHandle, frames, 0);
	if (frames < 0)
	{
		printf("Error playing wave: %s\n", snd_strerror(frames));
		return -1;
	}
	return frames;
}

int play_bytes(const void *bytes, int len) {
	snd_pcm_uframes_t frames, count;
	snd_pcm_uframes_t bufsize, period_size;
	frames = 0;
	count = 0;

	int mod = 4;
	if(len%mod != 0) {
		len += (mod - (len%mod));
	}

	snd_pcm_prepare(PlaybackHandle);

	snd_pcm_get_params(PlaybackHandle, &bufsize, &period_size);

	bufsize *= 2;
	printf("bufsize=%d\n", (int) bufsize);

	do {
		int remaining = len - count;
		int buflen = remaining < bufsize ? remaining : bufsize;
		frames = snd_pcm_writei(PlaybackHandle, bytes + count, buflen);
		// If an error, try to recover from it
		if (frames == -EPIPE) {
			printf("EPIPE\n");
			snd_pcm_prepare(PlaybackHandle);
		}
		if (frames < 0) {
			printf("Recovering\n");
			frames = snd_pcm_recover(PlaybackHandle, frames, 0);
		}
		if (frames < 0)
		{
			printf("Error playing wave: %s\n", snd_strerror(frames));
			break;
		}

		// Update our pointer
		count += frames;
		//printf("count=%d len=%d\n", (int)count, len);

	} while (count < len);

	// Wait for playback to completely finish
	if (count == len)
		snd_pcm_drain(PlaybackHandle);
	return 0;
}

int close_playback() {
	snd_pcm_close(PlaybackHandle);
	return 0;
}

int alsa_main(int argc, char **argv) {
	if(argc < 1) {
		printf("Usage: %s <WAV-file>\n", argv[0]);
		return -1;
	}

	int fd;
	unsigned long long len;

	fd = open(argv[1], O_RDONLY);
	// Find the length
	len = lseek(fd, 0, SEEK_END);

	// Skip the first 44 bytes (header)
	lseek(fd, 44, SEEK_SET);
	len -= 44;

	char *data = malloc(len);
	read(fd, data, len);

	int idx = init_playback("default", 44100, 2);
	play_bytes(idx, data, len);
	close_playback(idx);
	return 0;
}

int pulse_main() {
}
