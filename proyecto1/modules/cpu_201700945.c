#include <linux/init.h>
#include <linux/sched/signal.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/mm.h>
#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/sched.h>
#include <linux/version.h>
#include <linux/fs.h>

MODULE_LICENSE("GPL");
MODULE_AUTHOR("Bryan Paez");
MODULE_DESCRIPTION("Módulo del CPU");

#if LINUX_VERSION_CODE >= KERNEL_VERSION(5, 15, 0)
#define HAVE_PROC_OPS
#endif

static int mostrar_info_procesos(struct seq_file *archivo, void *v)
{
    struct task_struct *proceso;
    struct task_struct *proceso_hijo;
    struct list_head *lista;
    unsigned long uso_memoria;

    int primer_proceso = 0;

    int corriendo = 0;
    int durmiendo = 0;
    int zombie = 0;
    int detenido = 0;

    seq_printf(archivo, "{\n\"procesos\":[\n");
    for_each_process(proceso)
    {
        if (proceso->mm)
        {
            uso_memoria = get_mm_rss(proceso->mm) << PAGE_SHIFT;
        }
        else
        {
            uso_memoria = 0;
        }
        if (primer_proceso == 0)
        {
            seq_printf(archivo, "{");
            primer_proceso = 1;
        }
        else
        {
            seq_printf(archivo, ",{");
        }
        seq_printf(archivo, "\"pid\":%d,\n", proceso->pid);
        seq_printf(archivo, "\"nombre\":\"%s\",\n", proceso->comm);
        seq_printf(archivo, "\"usuario\": %d,\n", proceso->cred->uid);
        seq_printf(archivo, "\"estado\":%ld,\n", proceso->__state);
        int porcentaje = (((uso_memoria / (1024 * 1024))) * 100) / (15685);
        seq_printf(archivo, "\"ram\":%d,\n", porcentaje);

        seq_printf(archivo, "\"hijos\":[\n");
        int primer_hijo = 0;
        list_for_each(lista, &(proceso->children))
        {
            proceso_hijo = list_entry(lista, struct task_struct, sibling);
            if (primer_hijo != 0)
            {
                seq_printf(archivo, ",{");
                seq_printf(archivo, "\"pid\":%d,\n", proceso_hijo->pid);
                seq_printf(archivo, "\"nombre\":\"%s\"\n", proceso_hijo->comm);
                seq_printf(archivo, "}\n");
            }
            else
            {
                seq_printf(archivo, "{");
                seq_printf(archivo, "\"pid\":%d,\n", proceso_hijo->pid);
                seq_printf(archivo, "\"nombre\":\"%s\"\n", proceso_hijo->comm);
                seq_printf(archivo, "}\n");
                primer_hijo = 1;
            }
        }
        primer_hijo = 0;
        seq_printf(archivo, "\n]");

        if (proceso->__state == 0)
        {
            corriendo += 1;
        }
        else if (proceso->__state == 1)
        {
            durmiendo += 1;
        }
        else if (proceso->__state == 4)
        {
            zombie += 1;
        }
        else
        {
            detenido += 1;
        }
        seq_printf(archivo, "}\n");
    }
    seq_printf(archivo, "],\n");
    seq_printf(archivo, "\"corriendo\":%d,\n", corriendo);
    seq_printf(archivo, "\"durmiendo\":%d,\n", durmiendo);
    seq_printf(archivo, "\"zombie\":%d,\n", zombie);
    seq_printf(archivo, "\"detenido\":%d,\n", detenido);
    seq_printf(archivo, "\"total\":%d\n", corriendo + durmiendo + zombie + detenido);
    seq_printf(archivo, "}\n");
    return 0;
}

static int abrir_cpu(struct inode *inode, struct file *archivo)
{
    return single_open(archivo, mostrar_info_procesos, NULL);
}

#ifdef HAVE_PROC_OPS
static const struct proc_ops operaciones = {
    .proc_open = abrir_cpu,
    .proc_read = seq_read,
    .proc_lseek = seq_lseek,
    .proc_release = single_release,
};
#else
static const struct file_operations operaciones = {
    .owner = THIS_MODULE,
    .open = abrir_cpu,
    .read = seq_read,
    .llseek = seq_lseek,
    .release = single_release,
};
#endif

static int inicio(void)
{
    proc_create("cpu_201700945", 0, NULL, &operaciones);
    printk(KERN_INFO "Bryan Paez\n");

    return 0;
}

static void __exit finalizar(void)
{
    remove_proc_entry("cpu_201700945", NULL);
    printk(KERN_INFO "Segundo Semestre 2023\n");
}

module_init(inicio);
module_exit(finalizar);
