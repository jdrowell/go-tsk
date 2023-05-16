#include <tsk/libtsk.h>

int go_dirwalk_callback(TSK_FS_FILE *fs_file, char *path, void *data);

int dirwalk_callback(TSK_FS_FILE *fs_file, const char *path, void *data);
