#include "wrapper.h"

int dirwalk_callback(TSK_FS_FILE *fs_file, const char *path, void *data) {
    // Cast the data pointer back to its original type
    int *count_ptr = (int *)data;
    // Increment the count
    (*count_ptr)++;
    go_dirwalk_callback(fs_file, (char *)path, data);
    return 0;
}
