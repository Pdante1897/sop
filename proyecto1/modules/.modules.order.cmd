cmd_/home/bryan/Escritorio/sop/sop/proyecto1/modules/modules.order := {   echo /home/bryan/Escritorio/sop/sop/proyecto1/modules/cpu_201700945.ko;   echo /home/bryan/Escritorio/sop/sop/proyecto1/modules/ram_201700945.ko; :; } | awk '!x[$$0]++' - > /home/bryan/Escritorio/sop/sop/proyecto1/modules/modules.order
