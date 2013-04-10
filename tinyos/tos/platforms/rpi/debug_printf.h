#ifndef __DEBUG_PRINTF_H__
#define __DEBUG_PRINTF_H__

#include <stdio.h>

#define RPI_DEBUG 1
#define CC2520RPI_DEBUG 1
#define CC2520RPI_KERNEL_DRIVER_DEBUG 1
#define TUN_DEBUG 1
//#define IOMANAGER_DEBUG 1
//#define RPI_INTERRUPT_DEBUG 1
#define UART_DEBUG 1
#define EUI64_DEBUG 1


#define DBG(...)\
  do {\
    flockfile(stdout);\
    printf("%s:%d:\t", __FILE__, __LINE__); \
    printf(__VA_ARGS__);\
    funlockfile(stdout);\
  } while (0)


// fprintf(stderr, "%s:%d\t%s:\t", __FILE__, __LINE__, __func__);

#define ERROR(...)\
  do {\
    flockfile(stderr);\
    fprintf(stderr, "%s:%d\t", __FILE__, __LINE__); \
    fprintf(stderr, __VA_ARGS__);\
    funlockfile(stderr);\
  } while (0)


#ifdef RPI_DEBUG
#define RPI_PRINTF(...) DBG(__VA_ARGS__)
#else
#define RPI_PRINTF(...)
#endif

#ifdef CC2520RPI_DEBUG
#define RADIO_PRINTF(...) DBG(__VA_ARGS__)
#else
#define RADIO_PRINTF(...)
#endif

#ifdef TUN_DEBUG
#define TUN_PRINTF(...) DBG(__VA_ARGS__)
#else
#define TUN_PRINTF(...)
#endif

#ifdef IOMANAGER_DEBUG
#define IOMANAGER_PRINTF(...) DBG(__VA_ARGS__)
#else
#define IOMANAGER_PRINTF(...)
#endif

#ifdef RPI_INTERRUPT_DEBUG
#define INT_PRINTF(...) DBG(__VA_ARGS__)
#else
#define INT_PRINTF(...)
#endif

#ifdef UART_DEBUG
#define UART_PRINTF(...) DBG(__VA_ARGS__)
#else
#define UART_PRINTF(...)
#endif

#ifdef EUI64_DEBUG
#define EUI64_PRINTF(...) DBG(__VA_ARGS__)
#else
#define EUI64_PRINTF(...)
#endif

#endif
