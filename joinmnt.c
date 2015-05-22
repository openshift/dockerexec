#define _GNU_SOURCE
#include <stdlib.h>
#include <unistd.h>
#include <stdio.h>
#include <errno.h>
#include <string.h>

#include <linux/limits.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <sys/ioctl.h>
#include <fcntl.h>
#include <signal.h>
#include <sched.h>

#define pr_perror(fmt, ...) fprintf(stderr, "nsenter: " fmt ": %m\n", ##__VA_ARGS__)

// Use raw setns syscall for versions of glibc that don't include it (namely glibc-2.12)
#if __GLIBC__ == 2 && __GLIBC_MINOR__ < 14
#define _GNU_SOURCE
#include "syscall.h"
#ifdef SYS_setns
int setns(int fd, int nstype)
{
	return syscall(SYS_setns, fd, nstype);
}
#endif
#endif

void joinmnt()
{
	char *namespaces[] = { "mnt" };
	const int num = sizeof(namespaces) / sizeof(char *);
	char buf[PATH_MAX], *val;
	int i, tfd;
	pid_t docker_pid;

	val = getenv("_DOCKER_PID");
	if (val == NULL)
		return;

	docker_pid = atoi(val);
	snprintf(buf, sizeof(buf), "%d", docker_pid);
	if (strcmp(val, buf)) {
		pr_perror("Unable to parse _DOCKER_PID");
		exit(1);
	}

	/* Check that the specified process exists */
	snprintf(buf, PATH_MAX - 1, "/proc/%d/ns", docker_pid);
	tfd = open(buf, O_DIRECTORY | O_RDONLY);
	if (tfd == -1) {
		pr_perror("Failed to open \"%s\"", buf);
		exit(1);
	}

	for (i = 0; i < num; i++) {
		struct stat st;
		int fd;

		/* Symlinks on all namespaces exist for dead processes, but they can't be opened */
		if (fstatat(tfd, namespaces[i], &st, AT_SYMLINK_NOFOLLOW) == -1) {
			// Ignore nonexistent namespaces.
			if (errno == ENOENT)
				continue;
		}

		fd = openat(tfd, namespaces[i], O_RDONLY);
		if (fd == -1) {
			pr_perror("Failed to open ns file %s for ns %s", buf,
				  namespaces[i]);
			exit(1);
		}
		// Set the namespace.
		if (setns(fd, 0) == -1) {
			pr_perror("Failed to setns for %s", namespaces[i]);
			exit(1);
		}
		close(fd);
	}
}
