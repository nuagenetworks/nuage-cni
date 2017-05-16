import logging
import uuid
import time

import sys

from mesos.interface import Scheduler
from mesos.native import MesosSchedulerDriver
from mesos.interface import mesos_pb2

logging.basicConfig(level=logging.INFO)

agent_list = []
serviced_agent_list = []

def new_task(offer):
    task = mesos_pb2.TaskInfo()
    id = uuid.uuid4()
    task.task_id.value = str(id)
    task.slave_id.value = offer.slave_id.value
    task.name = "task {}".format(str(id))

    cpus = task.resources.add()
    cpus.name = "cpus"
    cpus.type = mesos_pb2.Value.SCALAR
    cpus.scalar.value = 1

    mem = task.resources.add()
    mem.name = "mem"
    mem.type = mesos_pb2.Value.SCALAR
    mem.scalar.value = 1

    return task


class InstallCNIScheduler(Scheduler):
    def __init__(self):
        self.master_ip = sys.argv[1]
        self.serviced_agents = 0

    def registered(self, driver, framework_id, master_info):
        logging.info("Registered with framework id: {}".format(framework_id))

    def resourceOffers(self, driver, offers):
        logging.info("Received resource offers: {}".format([o.id.value for o in offers]))
        for offer in offers:
            if offer.hostname not in agent_list:
                agent_list.append(offer.hostname)

        for offer in offers:
            if offer.hostname not in serviced_agent_list:
                serviced_agent_list.append(offer.hostname)
                task = new_task(offer)
                uri = task.command.uris.add()
	        # This framework will accept the url of the python Nuage CNI install script
                # hosted at a particular web server location as one of the command line args
	        uri.value = sys.argv[2] + '/install_nuage_cni.py'

	        # This script will also obtain the url of the web server hosting the 
                # components of Nuage CNI install as one of the command line args
                task.command.value = "python install_nuage_cni.py %s %s" % (sys.argv[2], sys.argv[3])
                time.sleep(2)
                logging.info("Launching task {task} "
                             "using offer {offer}.".format(task=task.task_id.value,
                                                       offer=offer.id.value))
                tasks = [task]
                driver.launchTasks(offer.id, tasks)

    def statusUpdate(self, driver, update):
        print "Task %s is in state %d" % (update.task_id.value, update.state)
        if update.state == mesos_pb2.TASK_FINISHED:
            print "Task successfully executed on agent node"
            self.serviced_agents += 1
            if self.serviced_agents == len(agent_list):
	        driver.stop()

if __name__ == '__main__':
    # Create a Mesos framework for CNI installation
    framework = mesos_pb2.FrameworkInfo()
    framework.user = ""  # Have Mesos fill in the current user.
    framework.name = "install-cni"
    driver = MesosSchedulerDriver(
        InstallCNIScheduler(),
        framework,
        sys.argv[1] + ":5050"  # assumes running on the master
    )
    driver.run()
