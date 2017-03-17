#ifndef __PULSE_H_
#define __PULSE_H_
int pa_init(const char *name, const char *device, int samplerate, int channels, int is_playback);
int pa_handle_read(int idx, const void *buf, int len);
int pa_handle_write(int idx, const void *buf, int len);
int pa_release(int idx);
int pa_get_latency(int idx);
#endif
