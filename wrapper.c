#include "wrapper.h"

int dirwalk_callback(TSK_FS_FILE *fs_file, const char *path, void *data) {
    go_dirwalk_callback(fs_file, (char *)path, data);
    return 0;
}
