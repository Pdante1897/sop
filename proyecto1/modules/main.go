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
    Memoria_total int `json:"memoria_total"`
    Memoria_libre int `json:"memoria_libre"`
	Memoria_en_uso      int     `json: "memoria_en_uso"`
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
    Usage     float64   `json: "usage"`
}

type CpuSend struct {
	Processes []ProcessSend `json: "processes"`
	Running   int           `json: "running"`
	Sleeping  int           `json: "sleeping"`
	Zombie    int           `json: "zombie"`
	Stopped   int           `json: "stopped"`
	Total     int           `json: "total"`
	Usage     float64       `json: "usage"`
}
type RAM struct {
	Used int `json: "memoria_en_uso"`
}

type ProcessSend struct {
	Pid   int      `json: "pid"`
	Name  string   `json: "name"`
	User  string   `json: "user"`
	State string   `json: "state"`
	Ram   int      `json: "ram"`
	Child []Childs `json: "child"`
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
    fmt.Println(string(ram))
	json.Unmarshal(ram, &memoria)
    fmt.Println(memoria)
	
	return memoria.Memoria_en_uso
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

	url := "http://35.245.67.156:4000/insertar_proceso" // Reemplaza con la URL correcta
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

	url := "http://35.245.67.156:4000/insertar_hijo" // Reemplaza con la URL correcta
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

	url := "http://35.245.67.156:4000/insertar_uso" // Reemplaza con la URL correcta
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

func InsertarTasks(maquina string, running string, sleeping string, zombie string, stopped string, total string) {
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

	url := fmt.Sprintf("http://35.245.67.156:4000/insertar_tarea/%s", maquina)
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
    memoria := string(getMemory())
	InsertarUsos(memoria, num1)
	runing := strconv.Itoa(cpu_info.Running)
	sleeping := strconv.Itoa(cpu_info.Sleeping)
	zombie  := strconv.Itoa(cpu_info.Zombie)
	stopped  := strconv.Itoa(cpu_info.Stopped)
	total  := strconv.Itoa(cpu_info.Total)
	InsertarTasks("1", runing, sleeping, zombie, stopped, total)
	for i:= 0; i < len(cpu_info.Processes); i++{
		Procesos := cpu_info.Processes[i]
		estado := getState(Procesos.State)
		pid := strconv.Itoa(Procesos.Pid)
		name := Procesos.Name 
		user := getUser(Procesos.User)
		ram := strconv.Itoa(Procesos.Ram) 
		if pid == `"`{
			fmt.Println("Se va alv >:v.")
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

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func killProcess(response http.ResponseWriter, request *http.Request) {
	data, errRead := ioutil.ReadAll(request.Body)
	fmt.Println("kill process")
	if errRead != nil {
		fmt.Println("error al eliminar un proceso")
		response.Write([]byte("{\"value\": false"))
	}
	fmt.Println(string(data))
	cmd := exec.Command("sh", "-c", `kill `+string(data))
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println("error al correr comando", err)
	}
	salida := strings.Trim(strings.Trim(string(stdout), " "), "\n")
	fmt.Println(salida)
	response.Write([]byte("{\"value\": true"))
}

func main() {
	
	tiempoDormido := (time.Second * time.Duration(5)) / time.Duration(1)
	for {
		LeecProcedimientos()
		time.Sleep(tiempoDormido)
	}
}