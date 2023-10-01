package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"log"
	"time"
	"io/ioutil"
	"encoding/json"
	"math/rand"
    "net/http"
    "bytes"
)
type Tree struct {
	Name  string   `json: "name"`
	Child []Childs `json: "child"`
}

type Memoria struct {
	Total_memory     int     `json: "total_memory"`
	Free_memory      int     `json: "free_memory"`
	Used_memory      int     `json: "used_memory"`
	Cache_memory     float64 `json: "cache_memory"`
	Available_memory int     `json: "available_memory"`
	MB_memory        int     `json: "mb_memory"`
}

type Childs struct {
	Pid  int    `json: "pid"`
	Name string `json: "name"`
}
type Process struct {
	Pid   int      `json: "pid"`
	Name  string   `json: "name"`
	User  int      `json: "user"`
	State int      `json: "state"`
	Ram   int      `json: "ram"`
	Child []Childs `json: "child"`
}


type CPU struct {
	Processes []Process `json: "processes"`
	Running   int       `json: "running"`
	Sleeping  int       `json: "sleeping"`
	Zombie    int       `json: "zombie"`
	Stopped   int       `json: "stopped"`
	Total     int       `json: "total"`
}

type CPU_SENDING struct {
	Running   int       `json: "running"`
	Sleeping  int       `json: "sleeping"`
	Zombie    int       `json: "zombie"`
	Stopped   int       `json: "stopped"`
	Total     int       `json: "total"`
}
type RAM struct {
	Used int `json: "used_memory"`
}

func getUser(nombre int) string {
	cmd := exec.Command("sh", "-c", `id -nu `+strconv.Itoa(nombre))
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println("error al correr comando", err)
	}
	salida := strings.Trim(strings.Trim(string(stdout), " "), "\n")
	//fmt.Println("Entro")
	//fmt.Println(salida)
	return salida
}
func getState(nombre int) string {
	if nombre == 0 {
		return "ejecucion"
	} else if nombre == 1 {
		return "dormido"
	} else if nombre == 4 {
		return "zombie"
	} else {
		return "detenido"
	}
}

func checkForErrors(err error) {
	if err != nil {
		panic(err.Error())
	}
}
func getCache() float64 {
	cmd := exec.Command("sh", "-c", `free -m | head -n2 | tail -1 | awk '{print $6}'`)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println("error al correr comando", err)
	}
	salida := strings.Trim(strings.Trim(string(stdout), " "), "\n")
	valor, _ := strconv.ParseFloat(salida, 64)
	return valor
}

func getMemory() int  {
	ram, _ := ioutil.ReadFile("/proc/ram_201700945")
	var memoria Memoria
	json.Unmarshal(ram, &memoria)
	memoria.Cache_memory = getCache()
	memoria.Used_memory = (memoria.Total_memory - memoria.Free_memory - int(getCache())) * 100 / memoria.Total_memory
	memoria.Available_memory = memoria.Free_memory + int(getCache())
	memoria.MB_memory = (memoria.Total_memory - memoria.Free_memory - int(getCache()))
	return memoria.Used_memory
}

func getCpuUsage() float64 {
	cmd := exec.Command("sh", "-c", `ps -eo pcpu | sort -k 1 -r | head -50`)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println("error al correr comando", err)
	}
	salidaAuxiliar := strings.Split(string(stdout), "\n")
	var total float64 = 0
	for i := 0; i < len(salidaAuxiliar); i++ {
		float1, _ := strconv.ParseFloat(salidaAuxiliar[i], 64)
		total += float1
	}
	total = (total / float64(len(salidaAuxiliar)-43))
	return (total)
}

func InsertarProceso(estado string, pid string, name string, user string, ram string) {
	data := map[string]string{
		"estado": estado,
		"pid":    pid,
		"name":   name,
		"user":   user,
		"ram":    ram,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	url := "http://192.168.32.2:4000/insertar_proceso" // Reemplaza con la URL correcta
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Proceso insertado correctamente.")
	} else {
		fmt.Println("Error al insertar proceso.")
	}
}

func InsertarTree(pidPadre string, pidHijo string, name string) {
	data := map[string]string{
		"pid_padre": pidPadre,
		"pid_hijo":  pidHijo,
		"name":      name,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	url := "http://192.168.32.2:4000/insertar_hijo" // Reemplaza con la URL correcta
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Hijo insertado correctamente.")
	} else {
		fmt.Println("Error al insertar hijo.")
	}
}

func InsertarUsos(ram string, cpu string) {
	data := map[string]string{
		"ram": ram,
		"cpu": cpu,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	url := "http://192.168.32.2:4000/insertar_uso" // Reemplaza con la URL correcta
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Uso insertado correctamente.")
	} else {
		fmt.Println("Error al insertar uso.")
	}
}

func InsertarTasks(running string, sleeping string, zombie string, stopped string, total string) {
	data := map[string]string{
		"running":  running,
		"sleeping": sleeping,
		"zombie":   zombie,
		"stopped":  stopped,
		"total":    total,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}

	url := "http://192.168.32.2:4000/insertar_tarea" // Reemplaza con la URL correcta
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Tarea insertada correctamente.")
	} else {
		fmt.Println("Error al insertar tarea.")
	}
}

func LeecProcedimientos(){
	Archivo, err := ioutil.ReadFile("/proc/cpu_201700945")
    if err != nil {
        log.Fatal(err)
    }
	cpu_info := CPU{}
	err = json.Unmarshal(Archivo, &cpu_info)
    if err != nil {
        log.Fatal(err)
    }
    rand.Seed(time.Now().UnixNano())
    num1 := strconv.Itoa(rand.Intn(30) + 20)
    num2 := strconv.Itoa(rand.Intn(30) + 20)
	InsertarUsos(num1, num2)
	runing := strconv.Itoa(cpu_info.Running)
	sleeping := strconv.Itoa(cpu_info.Sleeping)
	zombie  := strconv.Itoa(cpu_info.Zombie)
	stopped  := strconv.Itoa(cpu_info.Stopped)
	total  := strconv.Itoa(cpu_info.Total)
	InsertarTasks(runing, sleeping, zombie, stopped, total)
	for i:= 0; i < len(cpu_info.Processes); i++{
		Procesos := cpu_info.Processes[i]
		estado := getState(Procesos.State)
		pid := strconv.Itoa(Procesos.Pid)
		name := Procesos.Name 
		user := getUser(Procesos.User)
		ram := strconv.Itoa(Procesos.Ram) 
		if pid == `"`{
			fmt.Println("Se borra.")
		}else{
			InsertarProceso(estado,pid,name,user,ram)
		}
		fmt.Println(Procesos.Child)
		if len(Procesos.Child) > 0 {
			for j:= 0; j < len(Procesos.Child); j++{
				hijos := Procesos.Child[j]
				pid_hijo := strconv.Itoa(hijos.Pid)
				nombre := hijos.Name
				InsertarTree(pid, pid_hijo, nombre)
			}
		}else{
			fmt.Println("El padre no tiene hijos")
		}
	}
}
var DinamicTree []Tree 
var ProcessList []Process

func main() {
	
	tiempoDormido := (time.Second * time.Duration(5)) / time.Duration(1)
	for {
		LeecProcedimientos()
		time.Sleep(tiempoDormido)
	}
}