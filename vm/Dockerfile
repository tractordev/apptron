FROM i386/alpine:3.18.6

ENV KERNEL=lts
ENV HOSTNAME=localhost
ENV PASSWORD='root'

RUN apk add openrc \ 
            alpine-base \
            agetty \
            alpine-conf \
            github-cli \
            git \
            go \
            openssh

# Install mkinitfs from edge (todo: remove this when 3.19+ has worked properly with 9pfs)
RUN apk add mkinitfs --no-cache --allow-untrusted --repository https://dl-cdn.alpinelinux.org/alpine/edge/main/ 

RUN if [ "$KERNEL" == "lts" ]; then \
    apk add linux-lts \
            linux-firmware-none \
            linux-firmware-sb16; \
else \
    apk add linux-$KERNEL; \
fi

# Adding networking.sh script (works only on lts kernel yet)
RUN if [ "$KERNEL" == "lts" ]; then \ 
    echo -e "echo '127.0.0.1 localhost' >> /etc/hosts && rmmod ne2k-pci && modprobe ne2k-pci\nhwclock -s\nsetup-interfaces -a -r" > /root/networking.sh && \ 
    chmod +x /root/networking.sh; \ 
fi

RUN sed -i 's/getty 38400 tty1/agetty --autologin root tty1 linux/' /etc/inittab
RUN echo 'ttyS0::once:/sbin/agetty --autologin root -s ttyS0 115200 vt100' >> /etc/inittab 
RUN echo "root:$PASSWORD" | chpasswd

# assume env86 will add /bin/guest86
COPY ./guest86.service /etc/init.d/guest86
# there is a bug where it always crashes on first run
RUN echo "rc-service guest86 start && rc-service guest86 stop && rc-service guest86 start" > /root/.profile

# https://wiki.alpinelinux.org/wiki/Alpine_Linux_in_a_chroot#Preparing_init_services
RUN for i in devfs dmesg mdev hwdrivers; do rc-update add $i sysinit; done
RUN for i in hwclock modules sysctl hostname bootmisc; do rc-update add $i boot; done
RUN rc-update add killprocs shutdown

# Generate initramfs with 9p modules
RUN mkinitfs -F "ata base ide scsi virtio ext4 9p" $(cat /usr/share/kernel/$KERNEL/kernel.release)
