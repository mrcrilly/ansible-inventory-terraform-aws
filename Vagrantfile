Vagrant.configure("2") do |config|
  config.vm.box = "centos/7"
  config.vm.define "terraform-inventory"

  config.vm.provider :virtualbox do |v|
    v.memory = 1024
    v.cpus = 4
  end
  
  # This breaks due to the way CentOS have handled this internally
  config.vm.synced_folder ".", "/home/vagrant/sync", disabled: true

  config.vm.provision :shell, inline: "curl https://storage.googleapis.com/golang/go1.6.3.linux-amd64.tar.gz -o go.tgz"
  config.vm.provision :shell, inline: "tar -xzf go.tgz; sudo cp -r go /usr/local/"
end
