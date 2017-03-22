#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <errno.h>
#include <fcntl.h>
#include <pulse/simple.h>
#include <pulse/error.h>


struct config {
	pa_simple *pa_simple;
	pa_sample_spec spec;
	char *name;
	char *device;
	int is_playback;
};

struct config *configs[32];

int pa_init(char *name, char *device, int samplerate, int channels, int is_playback) {
	pa_simple *s;
	pa_sample_spec ss;

	ss.format = PA_SAMPLE_S16NE;
	ss.channels = channels;
	ss.rate = samplerate;

	pa_buffer_attr buffer_attrs = {128, 512, -1, -1};
	pa_buffer_attr *buf_attr_ptr = NULL;
	(void) buffer_attrs;
	/*
	buffer_attrs.maxlength = -1;
	buffer_attrs.tlength = -1;
	buffer_attrs.prebuf = -1;
	buffer_attrs.minreq = -1;
	buffer_attrs.fragsize = -1;
	*/

	int idx = 0;
	for(idx = 0; idx < 32; idx++) {
		if(configs[idx] == NULL) {
			break;
		}
	}
	if(idx == 32) {
		printf("Too many streams\n");
		return -1;
	}
	configs[idx] = malloc(sizeof(struct config));

	if(strcmp(device, "") == 0 || strcmp(device, "NULL") == 0) {
		device = NULL;
	}

	s = pa_simple_new(
			NULL,               // Use the default server.
			name,               // Our application's name.
			is_playback == 0 ? PA_STREAM_RECORD : PA_STREAM_PLAYBACK,
			device,
			//"<alsa_output.pci-0000_00_05.0.analog-stereo.monitor>",               // Use the default device.
			"libalsa pulse stream capture",            // Description of our stream.
			&ss,                // Our sample format.
			NULL,               // Use default channel map
			buf_attr_ptr,      // Use default buffering attributes.
			NULL                // Ignore error code.
			);
	if(s == NULL) {
		printf("Failed to open\n");
		return -1;
	}

	configs[idx]->pa_simple = s;
	configs[idx]->spec = ss;
	if(name != NULL) {
		configs[idx]->name = strdup(name);
	}
	if(device != NULL) {
		configs[idx]->device = strdup(device);
	}
	configs[idx]->is_playback = is_playback;
	return idx;
}

int pa_handle_read(int idx, char *buf, int len) {
	struct config *config = configs[idx];
	int error;

	/* Record some data ... */
	if (pa_simple_read(config->pa_simple, buf, len, &error) < 0) {
		fprintf(stderr, __FILE__": pa_simple_read() failed: %s\n", pa_strerror(error));
		return error;
	}
	return 0;
}

int pa_handle_write(int idx, char *buf, int len) {
	struct config *config = configs[idx];
	int error;

	/* Playback some data ... */
	if (pa_simple_write(config->pa_simple, buf, len, &error) < 0) {
		fprintf(stderr, __FILE__": pa_simple_read() failed: %s\n", pa_strerror(error));
		return error;
	}
	return 0;
}

int pa_get_latency(int idx) {
	struct config *config = configs[idx];
	int error;
	int latency;
	if((latency = pa_simple_get_latency(config->pa_simple, &error)) == -1) {
		printf("pa_simple_get_latency() failed: %s\n", pa_strerror(error));
		return -1;
	}
	return latency;
}

int pa_drain(int idx) {
	struct config *config = configs[idx];
	int error;
	if(pa_simple_drain(config->pa_simple, &error) < 0) {
		fprintf(stderr, __FILE__": pa_simple_drain() failed: %s\n", pa_strerror(error));
	}
	return error;
}

int pa_flush(int idx) {
	struct config *config = configs[idx];
	int error;
	if(pa_simple_flush(config->pa_simple, &error) < 0) {
		fprintf(stderr, __FILE__": pa_simple_flush() failed: %s\n", pa_strerror(error));
	}
	return error;
}

int pa_release(int idx) {
	struct config *config = configs[idx];
	if(config == NULL) {
		return -EINVAL;
	}
	printf("Releasing idx=%d\n", idx);
	free(config->name);
	free(config->device);
	pa_simple_free(config->pa_simple);
	printf("Released!\n");
	configs[idx] = NULL;
	return 0;
}
