package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"log"
	"main/structs"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	//"syscall"
	"time"

	//"github.com/gorilla/handlers"
	//"github.com/gorilla/mux"
	//"github.com/gorilla/websocket"
)



func main() {
	/*
	router := mux.NewRouter().StrictSlash(true)
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "PUT"})
	origins := handlers.AllowedOrigins([]string{"*"})
	port := os.Getenv("PORT")
	if port == "" {
		port = "4200"
	}
	router.HandleFunc("/", welcome).Methods("GET")
	router.HandleFunc("/ram", socketMemory)
	router.HandleFunc("/cpu", socketCpu)
	router.HandleFunc("/strace/{pid}", socketStrace)
	router.HandleFunc("/kill", killProcess).Methods("POST")
	router.HandleFunc("/cargarCpu", cargarCpu).Methods("GET")
	fmt.Println("servidor corriendo en puerto:" + port)
	http.ListenAndServe(":"+port, handlers.CORS(headers, methods, origins)(router))
	*/
    for {
        realizarPeticionGet("hola")
        time.Sleep(3000 * time.Millisecond)
    }
}

func realizarPeticionGet(ruta string) {
    url := os.Getenv("URL_API")
    if url == "" {
        url = "https://localhost:5000/"
    }

    respuesta, err := http.Get(url + ruta)
    if err != nil {
        fmt.Println("Error al realizar la solicitud GET:", err)
        return
    }
    defer respuesta.Body.Close()

    cuerpoRespuesta, err := ioutil.ReadAll(respuesta.Body)
    if err != nil {
        fmt.Println("Error al leer la respuesta:", err)
        return
    }

    fmt.Println(string(cuerpoRespuesta))
}

func realizarPeticionPost(ruta string) {
    url := os.Getenv("URL_API")
    if url == "" {
        url = "https://localhost:5000/"
    }

    datos := []byte(`{"clave": "valor"}`)

    respuesta, err := http.Post(url, "application/json", bytes.NewBuffer(datos))
    if err != nil {
        fmt.Println("Error al realizar la solicitud POST:", err)
        return
    }
    defer respuesta.Body.Close()

    cuerpoRespuesta, err := ioutil.ReadAll(respuesta.Body)

    fmt.Println(string(cuerpoRespuesta))
}

func cargarCpu(response http.ResponseWriter, request *http.Request) {
    numero := 123
    for i := 0; i < 100; i++ {
        numero = numero + numero
    }
    response.Write([]byte("exito"))
}

func obtenerUsoCPU() float64 {
    cmd := exec.Command("sh", "-c", `ps -eo pcpu | sort -k 1 -r | head -50`)
    stdout, err := cmd.Output()
    if err != nil {
        fmt.Println("error al ejecutar el comando", err)
    }
    salidaAuxiliar := strings.Split(string(stdout), "\n")
    var total float64 = 0
    for i := 0; i < len(salidaAuxiliar); i++ {
        float1, _ := strconv.ParseFloat(salidaAuxiliar[i], 64)
        total += float1
    }
    total = (total / float64(len(salidaAuxiliar) - 43))
    return total
}

func obtenerCache() float64 {
    cmd := exec.Command("sh", "-c", `free -m | head -n2 | tail -1 | awk '{print $6}'`)
    stdout, err := cmd.Output()
    if err != nil {
        fmt.Println("error al ejecutar el comando", err)
    }
    salida := strings.Trim(strings.Trim(string(stdout), " "), "\n")
    valor, _ := strconv.ParseFloat(salida, 64)
    return valor
}

func obtenerMemoria() structs.Memoria {
    ram, _ := ioutil.ReadFile("/proc/ram_mem_201700945")
    var memoria structs.Memoria
    json.Unmarshal(ram, &memoria)
    memoria.CacheMemory = obtenerCache()
    memoria.UsedMemory = (memoria.TotalMemory - memoria.FreeMemory - int(obtenerCache())) * 100 / memoria.TotalMemory
    memoria.AvailableMemory = memoria.FreeMemory + int(obtenerCache())
    memoria.MBMemory = (memoria.TotalMemory - memoria.FreeMemory - int(obtenerCache()))
    return memoria
}

func obtenerCPU() structs.CpuSend {
    processes, _ := ioutil.ReadFile("/proc/cpu_201700945")
    var cpu structs.Cpu
    var cpuSend structs.CpuSend
    json.Unmarshal(processes, &cpu)
    cpu.Usage = obtenerUsoCPU()
    hashmap := make(map[int]string)
    hashmap2 := make(map[int]string)
    var keys []int
    for i := 0; i < len(cpu.Processes); i++ {
        inputProcess := cpu.Processes[i]
        if !contiene(keys, inputProcess.User) {
            keys = append(keys, inputProcess.User)
            hashmap[inputProcess.User] = obtenerUsuario(inputProcess.User)
        }
        if !contiene(keys, inputProcess.State) {
            keys = append(keys, inputProcess.State)
            hashmap2[inputProcess.State] = obtenerEstado(inputProcess.State)
        }
        auxiliar := structs.ProcessSend{Pid: inputProcess.Pid, Name: inputProcess.Name, User: hashmap[inputProcess.User], State: hashmap2[inputProcess.State], Ram: inputProcess.Ram, Child: inputProcess.Child}
        cpuSend.Processes = append(cpuSend.Processes, auxiliar)
    }
    cpuSend.Running = cpu.Running
    cpuSend.Sleeping = cpu.Sleeping
    cpuSend.Stopped = cpu.Stopped
    cpuSend.Total = cpu.Total
    cpuSend.Usage = cpu.Usage
    cpuSend.Zombie = cpu.Zombie
    return cpuSend
}

func obtenerUsuario(nombre int) string {
    cmd := exec.Command("sh", "-c", `id -nu `+strconv.Itoa(nombre))
    stdout, err := cmd.Output()
    if err != nil {
        fmt.Println("error al ejecutar el comando", err)
    }
    salida := strings.Trim(strings.Trim(string(stdout), " "), "\n")
    return salida
}

func obtenerEstado(nombre int) string {
    switch nombre {
    case 0:
        return "ejecucion"
    case 1:
        return "dormido"
    case 4:
        return "zombie"
    default:
        return "detenido"
    }
}

func eliminarProceso(response http.ResponseWriter, request *http.Request) {
    data, errRead := ioutil.ReadAll(request.Body)
    fmt.Println("eliminar proceso")
    if errRead != nil {
        fmt.Println("error al eliminar un proceso")
        response.Write([]byte("{\"value\": false"))
    }
    fmt.Println(string(data))
    cmd := exec.Command("sh", "-c", `kill `+string(data))
    stdout, err := cmd.Output()
    if err != nil {
        fmt.Println("error al ejecutar el comando", err)
    }
    salida := strings.Trim(strings.Trim(string(stdout), " "), "\n")
    fmt.Println(salida)
    response.Write([]byte("{\"value\": true"))
}

func contiene(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}
