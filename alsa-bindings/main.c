#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <fcntl.h>
#include "alsa.h"
#include "pulse.h"

int alsa_main(char *data, int len) {
	int idx = init_playback("default", 44100, 2);
	printf("ALSA handle=%d\n", idx);
	play_bytes(idx, data, len);
	close_playback(idx);
	return 0;
}

int pulse_main(char *data, int len) {
	// Playback
	int idx = pa_init("playback", NULL, 48000, 2, 1);
	int bufsize = 24000;
	int count = 0;
	int remaining;
	int buflen;
	do {
		remaining = len - count;
		buflen = remaining < bufsize ? remaining : bufsize;
		pa_handle_write(idx, data + count, buflen);
		count += buflen;
	} while(remaining > 0);

	pa_release(idx);
	return 0;
}

int main(int argc, char **argv) {
	if(argc < 1) {
		printf("Usage: %s <WAV-file>\n", argv[0]);
		return -1;
	}

	int fd;
	unsigned int len = 0;

	fd = open(argv[1], O_RDONLY);

	struct file_header file_info;
	struct chunk_header chunk_info;
	struct format_header format_info;

	// Read in the file header
	read(fd, &file_info, sizeof(file_info));

	while(len == 0) {
		read(fd, &chunk_info, sizeof(chunk_info));
		if(strncmp(chunk_info.ID, "data", 4) == 0) {
			len = chunk_info.length;
			break;
		} else if(strncmp(chunk_info.ID, "fmt ", 4) == 0) {
			read(fd, &format_info, sizeof(format_info));
		}
	}
	printf("Length=%d\n", len);
	char *data = malloc(len);
	read(fd, data, len);

	//alsa_main(data, len);
	pulse_main(data, len);
	return 0;
}


