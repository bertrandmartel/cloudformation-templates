# CloudFormation templates

A collection of cloudformation templates

## Table of Content

* [VPC Infra](#vpc-infra)
* [ECS Infra](#ecs-infra)
* [ElastiCache](#elasticache)
* [Launch Lambda at creation](#lambda-launch-at-creation)

## VPC Infra

### 1 - Nat Instances

[nat-instances.yml](https://github.com/bertrandmartel/cloudformation-templates/blob/master/vpc-infra/nat-instances.yml)

![nat_instances](https://user-images.githubusercontent.com/5183022/73895563-52c6de80-4880-11ea-9678-7c8a29e93aa1.png)

Architecture featuring : 

* 1 VPC
* 3 public subnet accross 3AZ
* 3 private subnet accross 3AZ
* 3 NAT instances accross 3AZ
* 1 InternetGateway

The 3 NAT instances are living on different *public subnet* and have their own Security Group accepting incoming request (tcp & icmp) from their respective *private subnet*

Also 3 RoutingTable for each *private subnet* routes `0.0.0.0/0` to the corresponding NAT instance

### 2 - Nat Gateway

[nat-gateway.yml](https://github.com/bertrandmartel/cloudformation-templates/blob/master/vpc-infra/nat-gateway.yml)

![nat_gateway](https://user-images.githubusercontent.com/5183022/73895546-3a56c400-4880-11ea-8d14-15dd8a8aa81d.png)

Architecture featuring : 

* 1 VPC
* 3 public subnet accross 3AZ
* 3 private subnet accross 3AZ
* 3 NATGateway accross 3AZ
* 1 InternetGateway

3 RoutingTable for each *private subnet* routes `0.0.0.0/0` to the corresponding NAT Gateway

### Pricing

Note that using the Nat Gateway infra is 3 times more expensive than using Nat Instance (see [pricing](https://aws.amazon.com/vpc/pricing/))

## ECS Infra

[ecs-infra.yml](https://github.com/bertrandmartel/cloudformation-templates/blob/master/ecs/ecs-infra.yml)

[ecs-service.yml](https://github.com/bertrandmartel/cloudformation-templates/blob/master/ecs/ecs-service.yml)

![ecs-infra](https://user-images.githubusercontent.com/5183022/73895533-30cd5c00-4880-11ea-9824-e46d6b9917e4.png)

### ECS Infra stack 

Features : 

* 1 ECS Cluster
* 1 Application Load Balancer (ALB) on 3 Public subnet
* 1 Listener on 80 redirecting to https
* 1 Listener on 443 forwarding to a default Target Group
* 1 AutoScalingGroup on 3 Private subnet

The ALB Security Group accept incoming tcp requests from '0.0.0.0/0' on port 80 and 443

EC2 instances have SSM service installed (not installed by default on ECS optimized AMI)

### ECS Service stack

Features :

* 1 ECS Service
* 1 ECS TaskDefinition
* 1 Target Group which is targeted by the ECS Service
* 1 Listener Rule which route traffic for a specific hostname to the previous Target Group
* 1 Route53 DNS Record with the specific hostname pointing to the ALB (see ECS infra stack above)

Some notorious parameters :

* HostZone: name of your hostzone
* DockerImageURL: docker image URL

## ElastiCache

[elasticache.yml](https://github.com/bertrandmartel/cloudformation-templates/blob/master/elasticache/elasticache.yml)

* ElastiCache Cluster

Security Group of Elasticache Cluster accept incoming tcp request from '0.0.0.0/0' on port 6379

The subnet parameter should use private subnets

## Lambda launch at creation

In order to launch a lambda at stack creation, you need to use a [Custom Resource](https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-custom-resources.html). You need to create a lambda expecting to receive Custom Resource Request Input. Using go, it looks like [that](https://github.com/aws/aws-lambda-go/blob/master/cfn/event.go#L17-L26).

Also you need to create a physical resource ID when you receive the Create event. And then re-use the physical resource ID you will receive in the UPDATE or DELETE subsequent events. Failing to doing so would result in your stack getting stuck in `DELETE_IN_PROGRESS` state.

A minimal example in golang and the cloudformation stack :

* [lambda/main.go](https://github.com/bertrandmartel/cloudformation-templates/blob/master/lambda/main.go)
* [lambda/custom-resource.yml](https://github.com/bertrandmartel/cloudformation-templates/blob/master/lambda/custom-resource.yml)

Some useful resources about custom resources / lambda : 

* https://stories.schubergphilis.com/cloudformation-coding-with-custom-resources-9249c45bf37
* https://stackoverflow.com/a/41388529/2614364

