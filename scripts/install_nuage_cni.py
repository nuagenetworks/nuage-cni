import os
import sys

def install_nuage_cni():

    os.system("rpm -e nuage-cni")
    os.system("rm -irf /tmp/nuage-cni*.rpm")
    url = sys.argv[1] + '/' + sys.argv[2]
    os.system("wget %s -P /tmp" % url)
    os.system("rpm -i /tmp/nuage-cni*.rpm")

def main():

    install_nuage_cni()

if __name__ == "__main__":
    main()
