#ifndef __ALSA_H_
#define __ALSA_H_
struct file_header {
	char	ID[4];
	unsigned int	length;
	unsigned char	type[4];
};

struct chunk_header {
	char	ID[4];
	unsigned int	length;
};

struct format_header {
	short			format;
	unsigned short	channels;
	unsigned int	sample_rate;
	unsigned int	avg_bps;
	unsigned short	block_align;
	unsigned short	bits_per_sample;
};

int init_playback(const char *device, int samplerate, int channels);
int alsa_writei(const void *bytes, int len);
int play_bytes(int handle_idx, const void *bytes, int len);
int close_playback(int handle_idx);
#endif
