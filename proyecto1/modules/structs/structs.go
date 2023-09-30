package structs

type Prueba struct {
    Nombre string
    Edad   int
}

type Memoria struct {
    TotalMemory     int     `json:"total_memory"`
    FreeMemory      int     `json:"free_memory"`
    UsedMemory      int     `json:"used_memory"`
    CacheMemory     float64 `json:"cache_memory"`
    AvailableMemory int     `json:"available_memory"`
    MBMemory        int     `json:"mb_memory"`
}

type Process struct {
    Pid   int      `json:"pid"`
    Name  string   `json:"name"`
    User  int      `json:"user"`
    State int      `json:"state"`
    Ram   int      `json:"ram"`
    Child []Childs `json:"child"`
}

type Cpu struct {
    Processes []Process `json:"processes"`
    Running   int       `json:"running"`
    Sleeping  int       `json:"sleeping"`
    Zombie    int       `json:"zombie"`
    Stopped   int       `json:"stopped"`
    Total     int       `json:"total"`
    Usage     float64   `json:"usage"`
}

type ProcessSend struct {
    Pid   int      `json:"pid"`
    Name  string   `json:"name"`
    User  string   `json:"user"`
    State string   `json:"state"`
    Ram   int      `json:"ram"`
    Child []Childs `json:"child"`
}

type Childs struct {
    Pid  int    `json:"pid"`
    Name string `json:"name"`
}

type CpuSend struct {
    Processes []ProcessSend `json:"processes"`
    Running   int           `json:"running"`
    Sleeping  int           `json:"sleeping"`
    Zombie    int           `json:"zombie"`
    Stopped   int           `json:"stopped"`
    Total     int           `json:"total"`
    Usage     float64       `json:"usage"`
}

type StraceSend struct {
    Name string
    Pid  int
    List []Strace
}

type Strace struct {
    Recurrence int
    Name       string
}
