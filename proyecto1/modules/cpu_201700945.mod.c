#include <linux/module.h>
#define INCLUDE_VERMAGIC
#include <linux/build-salt.h>
#include <linux/elfnote-lto.h>
#include <linux/vermagic.h>
#include <linux/compiler.h>

BUILD_SALT;
BUILD_LTO_INFO;

MODULE_INFO(vermagic, VERMAGIC_STRING);
MODULE_INFO(name, KBUILD_MODNAME);

__visible struct module __this_module
__section(".gnu.linkonce.this_module") = {
	.name = KBUILD_MODNAME,
	.init = init_module,
#ifdef CONFIG_MODULE_UNLOAD
	.exit = cleanup_module,
#endif
	.arch = MODULE_ARCH_INIT,
};

#ifdef CONFIG_RETPOLINE
MODULE_INFO(retpoline, "Y");
#endif

static const struct modversion_info ____versions[]
__used __section("__versions") = {
	{ 0x4cf819e6, "module_layout" },
	{ 0x8e3e5524, "single_release" },
	{ 0xa5457562, "seq_lseek" },
	{ 0x83599202, "seq_read" },
	{ 0xa81d3524, "remove_proc_entry" },
	{ 0x92997ed8, "_printk" },
	{ 0xb4e4d6a9, "proc_create" },
	{ 0xe8f001e9, "init_task" },
	{ 0xa68165d, "seq_printf" },
	{ 0x5b8239ca, "__x86_return_thunk" },
	{ 0x801e8735, "single_open" },
	{ 0xbdfb6dbb, "__fentry__" },
};

MODULE_INFO(depends, "");


MODULE_INFO(srcversion, "445D88E1BF9E7A2B6729D6D");
