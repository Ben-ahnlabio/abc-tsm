# abc tsm stack

tsm 을 배포하기 위한 cloudformation stack

배포를 위해서 KMS, Core, Node Stack 이 필요하다.

https://docs.aws.amazon.com/ko_kr/secretsmanager/latest/userguide/auth-and-access_examples_cross.html

## 각 계정에서 해야할 일

### 배포 시

1. 각 계정에서 KMS stack 배포
2. Resource Share 초대 허용하기
3. ABC 계정에서는 EC2 AMI 를 대상 AWS 계정에 공유 설정함
4. 각 계정에서 Core Stack 배포
5. 필요한 경우 ABC 계정에서 Node Stack 배포

### Controller 업데이트 시

1. ABC 계정에서 업데이트된 AMI 를 대상 AWS 계정에 공유 설정
2. Core Stack 으로 업데이트
3. Node Stack 이 배포된 경우, Node Stack 업데이트


## KMS Stack

운영에 필요한 API Key, KMS Secret 등을 생성한다.

node1 과 node2 가 서로 다른 AWS 계정에 배포된다고 가정

abc-kms-stack.yaml 파일을 aws cloudformation 에 배포



### Parameters

#### Stack Name

dev 환경

- `abc-tsm-kms-node1-dev`
- `abc-tsm-kms-node2-dev`

production 환경

- `abc-tsm-kms-node1`
- `abc-tsm-kms-node2`


#### Namespace

dev 환경

- `abcdev`

production 환경

- `abc`

#### NodeAWSAccount

서로 다른 AWS 계정을 사용한다고 가정하며, node1 배포시 node2 가 배포되는 계정의 ID,
node2 배포시에는 node1 이 배포되는 계정의 ID 를 입력해야 한다.

#### DeletePolicy

Delete policy 는 cloudformation stack 을 삭제하는 경우, stack 에 의해 생성된 resource 를 함께 삭제할지 결정하는 값이다.

`Delete` 로 설정된 경우 Cloudformation Stack 을 삭제하면 리소스가 삭제된다.

`Retain` 으로 설정한 경우 Cloudformation Stack 이 생성한 리소스를 수작업으로 삭제해야 한다.


dev 환경

- `Delete` (기본값)

prod 환경

- `Retain`



## Node Stack

### Parameters

#### Stack Name

dev 환경

- `abc-tsm-core-node1-dev`
- `abc-tsm-core-node2-dev`

production 환경

- `abc-tsm-core-node1`
- `abc-tsm-core-node2`

vpc
subnets (priv, pub)
igw
nat gateway

ec2
security group
load balancer

route53

database