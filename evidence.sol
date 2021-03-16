// SPDX-License-Identifier: GPL-3.0
pragma solidity >=0.5.0 <0.6.0;

contract Evidence{
    uint SUCESS = 0;
    uint FILE_EXIST = 1;
    uint USER_NOT_EXIST = 2;
    uint FILEHASH_NOT_EXIST = 3;

    struct File {
        string hash;
        uint createTime;
        address owner;
    }
    
    // owner => (hash => File)
    mapping(address => mapping(string => File) ) evidenceMap;

    address[] userList;
    
    constructor () public {}
    
    function save (string memory fileHashHex) public returns (uint code,uint createTime){
        mapping(string => File) storage fileMap = evidenceMap[msg.sender];
        
        File storage f = fileMap[fileHashHex];
        if (f.createTime == 0){
            userList.push(msg.sender);
        }
        
        f.hash = fileHashHex;
        f.createTime = now;
        f.owner = msg.sender;
        
        evidenceMap[msg.sender][fileHashHex] = f;
        
        return (SUCESS,f.createTime);
    }
    
    function checkHash (string memory fileHashHex) public view returns(uint code) {
        if (evidenceMap[msg.sender][fileHashHex].createTime == 0){
            return (FILEHASH_NOT_EXIST);
        }
        
        return (FILE_EXIST);
    }
    
    function getEvidence (string memory fileHashHex) public view returns (uint code,uint createTime){
        File storage f = evidenceMap[msg.sender][fileHashHex];
        if (f.createTime == 0) {
            return (FILEHASH_NOT_EXIST,0);
        }
        
        return (SUCESS,f.createTime);
    }
    
    function getUsers () public view returns(address[] memory users){
        return userList;
    }
}