obj-m += cpu_201700945.o
obj-m += ram_201700945.o
all:
	make -C /lib/modules/$(shell uname -r)/build M=$(shell pwd) modules

clean:
	make -C /lib/modules/$(shell uname -r)/build M=$(shell pwd) clean
run:
	sudo insmod cpu_201700945.ko
	sudo insmod ram_201700945.ko 
