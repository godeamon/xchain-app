package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/contract"
	contractaccount "github.com/xuperchain/xuper-sdk-go/contract_account"
)

var (
	node   = "127.0.0.1:37101"
	bcname = "xuper"

	ce *contract.EVMContract

	mnemonic           = "聚 悬 带 肌 曹 术 时 别 浓 才 疾 保"
	contractAccount    = "XC9999999999999998@xuper"
	saveEvidenceMethod = "save"
	checkHashMethod    = "checkHash"
	getEvidenceMethod  = "getEvidence"
	getUsersMethod     = "getUsers"

	abi          = `[{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"constant":true,"inputs":[{"internalType":"string","name":"fileHashHex","type":"string"}],"name":"checkHash","outputs":[{"internalType":"uint256","name":"code","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"internalType":"string","name":"fileHashHex","type":"string"}],"name":"getEvidence","outputs":[{"internalType":"uint256","name":"code","type":"uint256"},{"internalType":"uint256","name":"createTime","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getUsers","outputs":[{"internalType":"address[]","name":"users","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"internalType":"string","name":"fileHashHex","type":"string"}],"name":"save","outputs":[{"internalType":"uint256","name":"code","type":"uint256"},{"internalType":"uint256","name":"createTime","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
	bin          = "60806040526000805560018055600280556003805534801561002057600080fd5b506109a1806100306000396000f3fe608060405234801561001057600080fd5b506004361061004b5760003560e01c8062ce8e3e1461005057806338e48f06146100af578063b16c6ee714610185578063e670f7cd1461025b575b600080fd5b61005861032a565b6040518080602001828103825283818151815260200191508051906020019060200280838360005b8381101561009b578082015181840152602081019050610080565b505050509050019250505060405180910390f35b610168600480360360208110156100c557600080fd5b81019080803590602001906401000000008111156100e257600080fd5b8201836020820111156100f457600080fd5b8035906020019184600183028401116401000000008311171561011657600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506103b8565b604051808381526020018281526020019250505060405180910390f35b61023e6004803603602081101561019b57600080fd5b81019080803590602001906401000000008111156101b857600080fd5b8201836020820111156101ca57600080fd5b803590602001918460018302840111640100000000831117156101ec57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610699565b604051808381526020018281526020019250505060405180910390f35b6103146004803603602081101561027157600080fd5b810190808035906020019064010000000081111561028e57600080fd5b8201836020820111156102a057600080fd5b803590602001918460018302840111640100000000831117156102c257600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610778565b6040518082815260200191505060405180910390f35b606060058054806020026020016040519081016040528092919081815260200182805480156103ae57602002820191906000526020600020905b8160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019060010190808311610364575b5050505050905090565b6000806000600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff1681526020019081526020016000209050600081856040518082805190602001908083835b602083106104355780518252602082019150602081019050602083039250610412565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902090506000816001015414156104de5760053390806001815401808255809150509060018203906000526020600020016000909192909190916101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff160217905550505b848160000190805190602001906104f6929190610840565b50428160010181905550338160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff16021790555080600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020866040518082805190602001908083835b602083106105b75780518252602082019150602081019050602083039250610594565b6001836020036101000a0380198251168184511680821785525050505050509050019150509081526020016040518091039020600082018160000190805460018160011615610100020316600290046106119291906108c0565b50600182015481600101556002820160009054906101000a900473ffffffffffffffffffffffffffffffffffffffff168160020160006101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055509050506000548160010154935093505050915091565b6000806000600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020846040518082805190602001908083835b6020831061071157805182526020820191506020810190506020830392506106ee565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902090506000816001015414156107655760035460008090509250925050610773565b600054816001015492509250505b915091565b600080600460003373ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff168152602001908152602001600020836040518082805190602001908083835b602083106107ee57805182526020820191506020810190506020830392506107cb565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902060010154141561083557600354905061083b565b60015490505b919050565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f1061088157805160ff19168380011785556108af565b828001600101855582156108af579182015b828111156108ae578251825591602001919060010190610893565b5b5090506108bc9190610947565b5090565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106108f95780548555610936565b8280016001018555821561093657600052602060002091601f016020900482015b8281111561093557825482559160010191906001019061091a565b5b5090506109439190610947565b5090565b61096991905b8082111561096557600081600090555060010161094d565b5090565b9056fea265627a7a72315820df60b56cc52d0f2db57231c429413ce141b5f81aa7a355a6de0fa13f9b7a958964736f6c63430005110032"
	contractName = "evidence33"

	evidenceFile = "./evidence_file"
)

func main() {
	acc, err := account.GetAccountFromPlainFile("./keys")
	if err != nil {
		panic(err)
	}

	// 创建合约账户：XC9999999999999998@xuper，如果已经创建好了不会重复创建。
	createContractAccount(acc)
	time.Sleep(time.Second * 5)

	// 部署合约。
	deployContract(acc)

	// 检查文件是否上链，此时还未上链。
	checkHash(acc)

	// 将文件哈希上链。
	save(acc)

	// 再次检查是否已经存储到链（如果文件未修改则成功，文件修改过则失败）。
	checkHash(acc)

	// 获取链上文件哈希数据。
	getEvidence(acc)

	// 更改文件内容
	modifyFile(evidenceFile)

	// 再次检查是否已经存储到链（如果文件未修改则成功，文件修改过则失败）。
	checkHash(acc)
}

func createAccount() {
	acc, err := account.CreateAccount(1, 1)
	if err != nil {
		panic(err)
	}
	log.Println(acc)
}

// createContractAccount 创建合约账户
func createContractAccount(acc *account.Account) {

	// 实例化一个可以创建合约账户的客户端对象
	ca := contractaccount.InitContractAccount(acc, node, bcname)

	// 发送创建合约账户交易
	// contractAccount = "XC9999999999999998@xuper"
	txid, err := ca.CreateContractAccount(contractAccount)
	if err != nil {
		log.Println("合约账户已经创建")
	} else {
		log.Println(txid)
	}
}

func deployContract(acc *account.Account) {
	// 创建合约操作客户端
	ec := getContractEVMClient(acc)

	// 部署合约，参数：合约初始化参数，合约 bin，合约 abi
	txid, err := ec.Deploy(nil, []byte(bin), []byte(abi))
	if err != nil {
		panic(err)
	}

	log.Println("部署合约成功，交易ID：", txid)
}

func save(acc *account.Account) {
	ec := getContractEVMClient(acc)

	// 存证合约的 save 方法的参数名字为 fileHashHex
	args := map[string]string{
		"fileHashHex": getFileHash(evidenceFile),
	}

	// saveEvidenceMethod = "save"
	txid, err := ec.Invoke(saveEvidenceMethod, args, "")
	if err != nil {
		panic(err)
	}
	log.Println("保存文件哈希成功，交易ID:", txid)
}

func getEvidence(acc *account.Account) {
	ec := getContractEVMClient(acc)

	args := map[string]string{
		// 存证合约的 getEvidence 方法的参数名字为 fileHashHex
		"fileHashHex": getFileHash(evidenceFile),
	}

	// getEvidenceMethod = "getEvidence"
	preExeRPCRes, err := ec.Query(getEvidenceMethod, args)
	if err != nil {
		panic(err)
	}

	log.Println("获取文件哈希交易结果：")
	for _, res := range preExeRPCRes.GetResponse().GetResponse() {
		printResp(res)
	}
}

func checkHash(acc *account.Account) {
	ec := getContractEVMClient(acc)

	args := map[string]string{
		"fileHashHex": getFileHash(evidenceFile),
	}
	preExeRPCRes, err := ec.Query(checkHashMethod, args)
	if err != nil {
		panic(err)
	}

	log.Println("检查文件哈希交易结果：")
	for _, res := range preExeRPCRes.GetResponse().GetResponse() {
		printResp(res)
	}
}

func printResp(res []byte) {

	result := []map[string]string{}
	err := json.Unmarshal(res, &result)
	if err != nil {
		fmt.Println("解析失败:", err)
		return
	}

	for _, v := range result {
		if code, ok := v["code"]; ok {
			switch code {
			case "0":
				log.Println("交易成功")
			case "1":
				log.Println("文件已存在")
			case "3":
				log.Println("文件哈希不存在")
			default:
				log.Println("invalid resp:", result)
			}
		}
		if c, ok := v["createTime"]; ok {
			log.Println("createTime:", c)
		}
	}
}

func getContractEVMClient(acc *account.Account) *contract.EVMContract {
	if ce != nil {
		return ce
	}

	return contract.InitEVMContract(acc, node, bcname, contractName, contractAccount)
}

func getFileHash(filePath string) string {
	h := sha256.New()
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := io.Copy(h, f); err != nil {
		panic(err)
	}

	return hex.EncodeToString(h.Sum(nil))
}

func modifyFile(filePath string) {
	log.Println("修改文件内容")
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString("a")
	write.Flush()
}
