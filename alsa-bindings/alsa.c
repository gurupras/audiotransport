#include <stdio.h>
#include <unistd.h>
#include <alsa/asoundlib.h>

#include "alsa.h"

struct config {
	snd_pcm_t *handle;
	char *device;
	int samplerate;
	int channels;
};

static struct config *configs[32];

int init_playback(const char *device, int samplerate, int channels) {
	int err;

	snd_pcm_t *PlaybackHandle;

	printf("Init parameters: %s %d %d\n", device, samplerate, channels);

	if((err = snd_pcm_open(&PlaybackHandle, device, SND_PCM_STREAM_PLAYBACK, 0)) < 0) {
		printf("Can't open audio %s: %s\n", device, snd_strerror(err));
		return -1;
	}

	if ((err = snd_pcm_set_params(PlaybackHandle, SND_PCM_FORMAT_S16, SND_PCM_ACCESS_RW_INTERLEAVED, channels, samplerate, 1, 500000)) < 0) {
		printf("Can't set sound parameters: %s\n", snd_strerror(err));
		return -1;
	}

	int i;
	for(i = 0; i < 32; i++) {
		if(configs[i] == NULL) {
			break;
		}
	}
	if(i == 32) {
		printf("Ran out of handles\n");
		return -1;
	}
	configs[i] = malloc(sizeof(struct config));

	struct config *config = configs[i];

	config->handle = PlaybackHandle;
	config->device = strdup(device);
	config->samplerate = samplerate;
	config->channels = channels;

	return i;
}

int alsa_writei(const void *bytes, int len) {
	return 0;
}

int play_bytes(int handle_idx, const void *bytes, int len) {
	snd_pcm_uframes_t frames, count;
	snd_pcm_uframes_t bufsize, period_size;
	frames = 0;
	count = 0;

	struct config *config = configs[handle_idx];

	snd_pcm_prepare(config->handle);
	snd_pcm_get_params(config->handle, &bufsize, &period_size);
	printf("bufsize=%d\n", (int) bufsize);

	do {
		int remaining = len - count;
		int buflen = remaining < bufsize ? remaining : bufsize;
		frames = snd_pcm_writei(config->handle, bytes + count, buflen);
		// If an error, try to recover from it
		if (frames == -EPIPE) {
			printf("EPIPE\n");
			snd_pcm_prepare(config->handle);
		}
		if (frames < 0) {
			printf("Recovering\n");
			frames = snd_pcm_recover(config->handle, frames, 0);
		}
		if (frames < 0)
		{
			printf("Error playing wave: %s\n", snd_strerror(frames));
			break;
		}

		// Update our pointer
		count += (frames * 2 * config->channels);
		//printf("count=%d len=%d\n", (int)count, len);

	} while (count < len);

	// Wait for playback to completely finish
	if (count == len)
		snd_pcm_drain(config->handle);
	return 0;
}

int close_playback(int handle_idx) {
	struct config *config = configs[handle_idx];
	if(config == NULL || config->handle == NULL) {
		return -EFAULT;
	}
	snd_pcm_close(config->handle);
	free(config->device);
	configs[handle_idx] = NULL;
	return 0;
}
