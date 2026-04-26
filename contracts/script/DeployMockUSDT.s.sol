// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "forge-std/Script.sol";
import "../src/MockUSDT.sol";

contract DeployMockUSDT is Script {
    function run() external {
        uint256 deployerPrivateKey = vm.envUint("DEPLOYER_PRIVATE_KEY");
        vm.startBroadcast(deployerPrivateKey);

        MockUSDT token = new MockUSDT();
        console.log("MockUSDT deployed at:", address(token));
        console.log("Owner:", token.owner());
        console.log("Total supply:", token.totalSupply());

        vm.stopBroadcast();
    }
}
