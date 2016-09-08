# snap collector plugin - procstat

## Collected Metrics

Namespace                                        | Description
-------------------------------------------------|------------------------------------------------------------
/staples/procstat/\*/fds                          | The number of file descriptors in use by this process
/staples/procstat/\*/involuntary_context_switches | The number of involuntary context switches forced by the os on the process
/staples/procstat/\*/memory_rss                   | Physical memory in use by the process
/staples/procstat/\*/memory_swap                  | Swap memory in use by the process
/staples/procstat/\*/memory_vms                   | Virtual memory in use by the process
/staples/procstat/\*/numThreads                   | Number of active threads owned by the process
/staples/procstat/\*/read_bytes                   | Number of bytes this process caused to be fetched from the storage layer
/staples/procstat/\*/read_count                   | Number of read i/o operations by the process
/staples/procstat/\*/voluntary_context_switches   |
/staples/procstat/\*/write_bytes                  | Number of bytes this process caused to be written to the storage layer
/staples/procstat/\*/write_count                  | Number of write i/o operations by the process
/staples/procstat/\*/process_uptime               | Number of seconds the process has been alive
/staples/procstat/\*/cpu_time_guest               | The amount of time servicing guest OS systems by the process
/staples/procstat/\*/cpu_time_guest_nice          | The amount of time servicing guest OS systems by the process
/staples/procstat/\*/cpu_time_idle                | The amount of time the process spent idle
/staples/procstat/\*/cpu_time_iowait              | The amount of time spent by the process servicing servicing i/o waits
/staples/procstat/\*/cpu_time_irq                 | The amount of time servicing interrupts due to the process
/staples/procstat/\*/cpu_time_nice                | The amount of time the process spent in user mode at low priority
/staples/procstat/\*/cpu_time_soft_irq            | The amount of time servicing soft interrupts due to the process
/staples/procstat/\*/cpu_time_stolen              | The amount of stolen time for this process
/staples/procstat/\*/cpu_time_system              | The amount of time spent in system mode on the cpu for this process
/staples/procstat/\*/cpu_time_user                | The amount of time spent in user mode on the cpu for this process
