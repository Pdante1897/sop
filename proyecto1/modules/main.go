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
	Nombre  string   `json: "nombre"`
	Usuario  int      `json: "usuario"`
	Estado int      `json: "estado"`
	Ram   int      `json: "ram"`
}

type Procesos struct {
    Procesos []Process `json: "procesos"`
}
type ProcesosSend struct {
    Procesos []ProcessSend `json: "procesos"`
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
	Pid   string      `json: "pid"`
	Nombre  string   `json: "nombre"`
	Usuario  string      `json: "usuario"`
	Estado string      `json: "estado"`
	Ram   string      `json: "ram"`
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

func InsertarProceso(listado Procesos, maquina string, numero int) {
	//convertir listado de Procesos a listado de ProcessSend
    var listadoSend []ProcessSend
    for i:= 0; i < len(listado.Procesos); i++{
        estado := getState(listado.Procesos[i].Estado)
		pid := strconv.Itoa(listado.Procesos[i].Pid)
		name := listado.Procesos[i].Nombre 
		user := getUser(listado.Procesos[i].Usuario)
		ram := strconv.Itoa(listado.Procesos[i].Ram)
        listadoSend = append(listadoSend, ProcessSend{Pid: pid, Nombre: name, Usuario: user, Estado: estado, Ram: ram})
        //fmt.Println(user)
        //fmt.Println(listadoSend[i])
    }
    procesosSend := ProcesosSend{Procesos: listadoSend}
    data := procesosSend
    
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("http://35.245.67.156:4000/insertar_proceso/%s/%s", maquina, strconv.Itoa(numero))
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

func InsertarUsos(maquina string, ram string, cpu string) {
	data := map[string]string{
		"ram": ram,
		"cpu": cpu,
	}
	var respuesta struct {
        Message string                 `json:"message"`
        Data    map[string]interface{} `json:"data"`
    }
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("http://35.245.67.156:4000/insertar_uso/%s", maquina)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	cuerpo, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	fmt.Println(string(cuerpo))

	if err := json.Unmarshal([]byte(cuerpo), &respuesta); err != nil {
        fmt.Println("Error al deserializar JSON:", err)
	}

	pid := respuesta.Data["pid"]
    kill := respuesta.Data["kill"]

	fmt.Println(pid)
	fmt.Println(kill)

	if kill == "true" {
		killProcess(pid.(string))
	}

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

func LeecProcedimientos(iteracion int){
    Archivo, err := ioutil.ReadFile("/proc/cpu_201700945")
    if err != nil {
        log.Fatal(err)
    }
    //fmt.Println(string(Archivo))
	cpu_info := CPU{}
    procesos := Procesos{}
    err = json.Unmarshal(Archivo, &procesos)
    if err != nil {
        fmt.Println("Error al parsear el JSON:", err)
        return
    }
	err = json.Unmarshal(Archivo, &cpu_info)
    if err != nil {
        log.Fatal(err)
    }
    rand.Seed(time.Now().UnixNano())
    num1 := strconv.FormatFloat(getCpuUsage(), 'f', 2, 64)
    memoria := strconv.Itoa(getMemory())
    fmt.Println(memoria)
    //fmt.Println(procesos)

	InsertarUsos("1", memoria, num1)
    InsertarProceso(procesos, "1", iteracion)

	//InsertarTasks("1", runing, sleeping, zombie, stopped, total)
	/*
    for i:= 0; i < len(procesos.Procesos); i++{
		estado := getState(procesos.Procesos[i].Estado)
		pid := strconv.Itoa(procesos.Procesos[i].Pid)
		name := procesos.Procesos[i].Nombre 
		user := getUser(procesos.Procesos[i].Usuario)
		ram := strconv.Itoa(procesos.Procesos[i].Ram)
        fmt.Println(user) 
        fmt.Println(procesos.Procesos[i]) 
		
		
	}*/
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

func killProcess(data string) {
	fmt.Println("kill process")
	fmt.Println(data)
	cmd := exec.Command("sh", "-c", `kill `+data)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println("error al correr comando", err)
	}
	salida := strings.Trim(strings.Trim(string(stdout), " "), "\n")
	fmt.Println(salida)
}

func main() {
	
	tiempoDormido := (time.Second * time.Duration(5)) / time.Duration(1)
	i:=0
    for {
        
		LeecProcedimientos(i)
		time.Sleep(tiempoDormido)
        i++
    }
}