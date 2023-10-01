#include <linux/init.h>
#include <linux/kernel.h>
#include <linux/proc_fs.h>
#include <linux/sched.h>
#include <linux/seq_file.h>
#include <linux/sched/signal.h>
#include <linux/fs.h>
#include <linux/version.h>
#include <linux/mm.h>
#include <linux/module.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Bryan Paez");
MODULE_DESCRIPTION("Módulo de Información de Memoria RAM");

#if LINUX_VERSION_CODE >= KERNEL_VERSION(5, 15, 0)
#define HAVE_PROC_OPS
#endif

struct sysinfo informacion;

static int mostrar_info_memoria(struct seq_file *archivo, void *v)
{
    si_meminfo(&informacion);
    unsigned long memoria_total = (informacion.totalram * 4);
    unsigned long memoria_libre = (informacion.freeram * 4) - (informacion.sharedram * 4) - (informacion.bufferram * 4);
    seq_printf(archivo, "{\n");
    seq_printf(archivo, "\"memoria_total\": %lu,\n", memoria_total / 1024);
    seq_printf(archivo, "\"memoria_libre\": %lu,\n", memoria_libre / 1024);
    seq_printf(archivo, "\"memoria_en_uso\": %lu\n", ((memoria_total - memoria_libre) * 100) / memoria_total);
    seq_printf(archivo, "}\n");
    return 0;
}

static int abrir_memoria(struct inode *inode, struct file *archivo)
{
    return single_open(archivo, mostrar_info_memoria, NULL);
}

#ifdef HAVE_PROC_OPS
static const struct proc_ops operaciones = {
    .proc_open = abrir_memoria,
    .proc_read = seq_read,
    .proc_lseek = seq_lseek,
    .proc_release = single_release,
};
#else
static const struct file_operations operaciones = {
    .owner = THIS_MODULE,
    .open = abrir_memoria,
    .read = seq_read,
    .llseek = seq_lseek,
    .release = single_release,
};
#endif

static int inicio(void)
{
    proc_create("ram_201700945", 0, NULL, &operaciones);
    printk(KERN_INFO "Bryan paez\n");
    return 0;
}

static void __exit finalizar(void)
{
    remove_proc_entry("ram_201700945", NULL);
    printk(KERN_INFO "Segundo Semestre 2023\n");
}

module_init(inicio);
module_exit(finalizar);
