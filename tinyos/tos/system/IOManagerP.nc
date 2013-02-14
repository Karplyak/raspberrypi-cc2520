#include <stdio.h>
#include <stdlib.h>
#include <sys/time.h>
#include <sys/types.h>
#include <unistd.h>
#include <errno.h>

module IOManagerP {
  provides {
    interface IO[uint8_t id];
    interface BlockingIO;
  }
}

implementation {

  enum {
    N_FDS = uniqueCount("FILEID"),
  };

  uint8_t map[N_FDS];
  fd_set  rfds;
  bool is_init = FALSE;


  command error_t IO.registerFileDescriptor[uint8_t id] (int file_descriptor) {
    if (!is_init) {
      memset(map, 0x01, sizeof(uint8_t) * N_FDS);
      is_init = TRUE;
    }
    map[id] = file_descriptor;
    return SUCCESS;
  }

  async command void BlockingIO.waitForIO () {
    int     ret;
    int     i;
    uint8_t nfds = 0;

    FD_ZERO(&rfds);
    for (i=0; i<N_FDS; i++) {
      if (map[i] != 1) {
        FD_SET(map[i], &rfds);
        if (map[i] + 1 > nfds) {
          nfds = map[i] + 1;
        }
      }
    }

    ret = select(nfds, &rfds, NULL, NULL, NULL);

    if (ret < 0) {
      if (errno == EINTR) {
        // suppress
        printf("select intre\n");
      } else {
        // error
        fprintf(stderr, "IOManagerP: select return error: %i\n", ret);
      }

    } else if (ret == 0) {
      //
      fprintf(stderr, "IOManagerP: select return 0.\n");

    } else {
      // some file is ready
      int j;
      printf("file ready\n");
      for (j=0; j<N_FDS; j++) {
        if (FD_ISSET(map[j], &rfds)) {
          signal IO.receiveReady[j]();
        }
      }
    }
  }

  default async event void IO.receiveReady[uint8_t id] () { }

}
