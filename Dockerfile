FROM centos

# Run the file to lauch Nuage CNI daemon
# as a Docker container
CMD /usr/bin/nuage -daemon
