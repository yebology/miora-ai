// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

/**
 * @title MockUSDT
 * @notice Mock USDT token for Miora AI bot trading on Base Sepolia.
 * @dev Owner can mint tokens to any address. Used as bot budget denomination.
 *      Users deposit MockUSDT to the Agentic Wallet, bot trades with it.
 */
contract MockUSDT is ERC20, Ownable {
    constructor() ERC20("Mock USDT", "mUSDT") Ownable(msg.sender) {
        // Mint 1,000,000 mUSDT to deployer
        _mint(msg.sender, 1_000_000 * 10 ** decimals());
    }

    /// @notice Mint tokens to any address. Only owner.
    function mint(address to, uint256 amount) external onlyOwner {
        _mint(to, amount);
    }

    /// @notice Override decimals to 6 (same as real USDT).
    function decimals() public pure override returns (uint8) {
        return 6;
    }
}
