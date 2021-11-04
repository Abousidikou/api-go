-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema monitorDB
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema monitorDB
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `monitorDB` ;
USE `monitorDB` ;

-- -----------------------------------------------------
-- Table `monitorDB`.`Provider`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `monitorDB`.`Provider` (
  `Provider_id` INT NOT NULL AUTO_INCREMENT,
  `Provider_ISP` VARCHAR(255) NULL,
  `Provider_ORG` VARCHAR(255) NULL,
  `Provider_AS_Number` VARCHAR(45) NULL,
  `Provider_AS_Name` VARCHAR(255) NULL,
  PRIMARY KEY (`Provider_id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `monitorDB`.`Country`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `monitorDB`.`Country` (
  `Country_id` INT NOT NULL AUTO_INCREMENT,
  `Country_Name` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`Country_id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `monitorDB`.`Region`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `monitorDB`.`Region` (
  `Region_id` INT NOT NULL AUTO_INCREMENT,
  `Region_Name` VARCHAR(45) NOT NULL,
  `Region_Country_id` INT NULL,
  PRIMARY KEY (`Region_id`),
  INDEX `fk_Region_Country_idx` (`Region_Country_id` ASC) ,
  CONSTRAINT `fk_Region_Country`
    FOREIGN KEY (`Region_Country_id`)
    REFERENCES `monitorDB`.`Country` (`Country_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `monitorDB`.`City`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `monitorDB`.`City` (
  `City_id` INT NOT NULL AUTO_INCREMENT,
  `City_Name` VARCHAR(45) NOT NULL,
  `City_Region_id` INT NULL,
  PRIMARY KEY (`City_id`),
  INDEX `fk_City_Region_idx` (`City_Region_id` ASC) ,
  CONSTRAINT `fk_City_Region`
    FOREIGN KEY (`City_Region_id`)
    REFERENCES `monitorDB`.`Region` (`Region_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `monitorDB`.`BBRInfo`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `monitorDB`.`BBRInfo` (
  `BBRInfo_id` INT NOT NULL AUTO_INCREMENT,
  `AvgBW` INT NULL,
  `MinBW` INT NULL,
  `MaxBw` INT NULL,
  `MedianBW` INT NULL,
  `AvgMinRTT` INT NULL,
  `MinMinRTT` INT NULL,
  `MaxMinRTT` INT NULL,
  `MedianMinRTT` INT NULL,
  PRIMARY KEY (`BBRInfo_id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `monitorDB`.`TCPInfo`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `monitorDB`.`TCPInfo` (
  `TCPInfo_id` INT NOT NULL AUTO_INCREMENT,
  `AvgRetransmits` INT NULL,
  `MinRetransmits` INT NULL,
  `MaxRetransmits` INT NULL,
  `MedianRetransmits` INT NULL,
  `AvgBackoff` INT NULL,
  `MinBackoff` INT NULL,
  `MaxBackoff` INT NULL,
  `MedianBackoff` INT NULL,
  `AvgRTO` INT NULL,
  `MinRTO` INT NULL,
  `MaxRTO` INT NULL,
  `MedianRTO` INT NULL,
  `AvgATO` INT NULL,
  `MaxATO` INT NULL,
  `MinATO` INT NULL,
  `MedianATO` INT NULL,
  `AvgSndMSS` INT NULL,
  `MinSndMSS` INT NULL,
  `MaxSndMSS` INT NULL,
  `MedianSndMSS` INT NULL,
  `AvgRcvMSS` INT NULL,
  `MinRcvMSS` INT NULL,
  `MaxRcvMSS` INT NULL,
  `MedianRcvMSS` INT NULL,
  `AvgUnacked` INT NULL,
  `MinUnacked` INT NULL,
  `MaxUnacked` INT NULL,
  `MedianUnacked` INT NULL,
  `AvgSacked` INT NULL,
  `MinSacked` INT NULL,
  `MaxSacked` INT NULL,
  `MedianSacked` INT NULL,
  `AvgLost` INT NULL,
  `MinLost` INT NULL,
  `MaxLost` INT NULL,
  `MedianLost` INT NULL,
  `AvgRetrans` INT NULL,
  `MinRetrans` INT NULL,
  `MaxRetrans` INT NULL,
  `MedianRetrans` INT NULL,
  `AvgPMTU` INT NULL,
  `MinPMTU` INT NULL,
  `MaxPMTU` INT NULL,
  `MedianPMTU` INT NULL,
  `AvgRcvSsThresh` INT NULL,
  `MinRcvSsThresh` INT NULL,
  `MaxRcvSsThresh` INT NULL,
  `MedianRcvSsThresh` INT NULL,
  `AvgRTT` INT NULL,
  `MinRTT` INT NULL,
  `MaxRTT` INT NULL,
  `MedianRTT` INT NULL,
  `AvgRTTVar` INT NULL,
  `MinRTTVar` INT NULL,
  `MaxRTTVar` INT NULL,
  `MedianRTTVar` INT NULL,
  `AvgAdvMSS` INT NULL,
  `MinAdvMSS` INT NULL,
  `MaxAdvMSS` INT NULL,
  `MedianAdvMSS` INT NULL,
  `AvgReordering` INT NULL,
  `MinReordering` INT NULL,
  `MaxReordering` INT NULL,
  `MedianReordering` INT NULL,
  `AvgRcvSpace` INT NULL,
  `MinRcvSpace` INT NULL,
  `MaxRcvSpace` INT NULL,
  `MedianRcvSpace` INT NULL,
  `AvgTotalRetrans` INT NULL,
  `MinTotalRetrans` INT NULL,
  `MedianTotalRetrans` INT NULL,
  `AvgBytesAcked` INT NULL,
  `MinBytesAcked` INT NULL,
  `MaxBytesAcked` INT NULL,
  `MedianBytesAcked` INT NULL,
  `AvgBytesReceived` INT NULL,
  `MinBytesReceived` INT NULL,
  `MaxBytesReceived` INT NULL,
  `MedianBytesReceived` INT NULL,
  `AvgSegsOut` INT NULL,
  `MinSegsOut` INT NULL,
  `MaxSegsOut` INT NULL,
  `MedianSegsOut` INT NULL,
  `AvgSegsIn` INT NULL,
  `MinSegsIn` INT NULL,
  `MaxSegsIn` INT NULL,
  `MedianSegsIn` INT NULL,
  `AvgNotsentBytes` INT NULL,
  `MinNotsentBytes` INT NULL,
  `MaxNotsentBytes` INT NULL,
  `MedianNotsentBytes` INT NULL,
  `AvgDataSegsIn` INT NULL,
  `MinDataSegsIn` INT NULL,
  `MaxDataSegsIn` INT NULL,
  `MedianDataSegsIn` INT NULL,
  `AvgDataSegsOut` INT NULL,
  `MinDataSegsOut` INT NULL,
  `MaxDataSegsOut` INT NULL,
  `MedianDataSegsOut` INT NULL,
  `AvgDeliveryRate` INT NULL,
  `MinDeliveryRate` INT NULL,
  `MaxDeliveryRate` INT NULL,
  `MedianDeliveryRate` INT NULL,
  `AvgBusyTime` INT NULL,
  `MinBusyTime` INT NULL,
  `MaxBusyTime` INT NULL,
  `MedianBusyTime` INT NULL,
  `AvgRWndLimited` INT NULL,
  `MinRWndLimited` INT NULL,
  `MaxRWndLimited` INT NULL,
  `MedianRWndLimited` INT NULL,
  `AvgSndBufLimited` INT NULL,
  `MinSndBufLimited` INT NULL,
  `MaxSndBufLimited` INT NULL,
  `MedianSndBufLimited` INT NULL,
  `AvgDelivered` INT NULL,
  `MinDelivered` INT NULL,
  `MaxDelivered` INT NULL,
  `MedianDelivered` INT NULL,
  `AvgDeliveredCE` INT NULL,
  `MinDeliveredCE` INT NULL,
  `MaxDeliveredCE` INT NULL,
  `MedianDeliveredCE` INT NULL,
  `AvgBytesSent` INT NULL,
  `MinBytesSent` INT NULL,
  `MaxBytesSent` INT NULL,
  `MedianBytesSent` INT NULL,
  `AvgBytesRetrans` INT NULL,
  `MinBytesRetrans` INT NULL,
  `MaxBytesRetrans` INT NULL,
  `MedianBytesRetrans` INT NULL,
  `AvgDsackDups` INT NULL,
  `MinDsackDups` INT NULL,
  `MaxDsackDups` INT NULL,
  `MedianDsackDups` INT NULL,
  `AvgReordSeen` INT NULL,
  `MinReordSeen` INT NULL,
  `MaxReordSeen` INT NULL,
  `MedianReordSeen` INT NULL,
  `MaxTotalRetrans` INT NULL,
  PRIMARY KEY (`TCPInfo_id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `monitorDB`.`Service`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `monitorDB`.`Service` (
  `Service_id` INT NOT NULL AUTO_INCREMENT,
  `Service_level` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`Service_id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `monitorDB`.`DaySlice`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `monitorDB`.`DaySlice` (
  `DaySlice_id` INT NOT NULL AUTO_INCREMENT,
  `DaySlice_level` VARCHAR(45) NOT NULL,
  PRIMARY KEY (`DaySlice_id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `monitorDB`.`Tests`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `monitorDB`.`Tests` (
  `Test_id` INT NOT NULL AUTO_INCREMENT,
  `Test_UUID` VARCHAR(255) NULL,
  `Test_Type` VARCHAR(10) NULL,
  `Test_ServerIP` VARCHAR(255) NULL,
  `Test_ServerPort` INT NULL,
  `Test_ClientIP` VARCHAR(255) NULL,
  `Test_ClientPort` INT NULL,
  `Test_Date` DATE NOT NULL,
  `Test_Country_id` INT NULL,
  `Test_Region_id` INT NULL,
  `Test_City_id` INT NULL,
  `Test_Provider_id` INT NULL,
  `Test_BBRInfo_id` INT NULL,
  `Test_TCPInfo_id` INT NULL,
  `Test_Service_id` INT NULL,
  `Test_DaySlice_id` INT NULL,
  PRIMARY KEY (`Test_id`),
  INDEX `fk_Tests_Provider_idx` (`Test_Provider_id` ASC) ,
  INDEX `fk_Tests_Country_idx` (`Test_Country_id` ASC) ,
  INDEX `fk_Tests_Region_idx` (`Test_Region_id` ASC) ,
  INDEX `fk_Tests_City_idx` (`Test_City_id` ASC) ,
  INDEX `fk_Tests_BBRInfo_idx` (`Test_BBRInfo_id` ASC) ,
  INDEX `fk_Tests_TCPInfo_idx` (`Test_TCPInfo_id` ASC) ,
  INDEX `fk_Tests_Service_idx` (`Test_Service_id` ASC) ,
  INDEX `fk_Tests_DaySlice_idx` (`Test_DaySlice_id` ASC) ,
  CONSTRAINT `fk_Tests_Provider`
    FOREIGN KEY (`Test_Provider_id`)
    REFERENCES `monitorDB`.`Provider` (`Provider_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_Tests_Country`
    FOREIGN KEY (`Test_Country_id`)
    REFERENCES `monitorDB`.`Country` (`Country_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_Tests_Region`
    FOREIGN KEY (`Test_Region_id`)
    REFERENCES `monitorDB`.`Region` (`Region_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_Tests_City`
    FOREIGN KEY (`Test_City_id`)
    REFERENCES `monitorDB`.`City` (`City_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_Tests_BBRInfo`
    FOREIGN KEY (`Test_BBRInfo_id`)
    REFERENCES `monitorDB`.`BBRInfo` (`BBRInfo_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_Tests_TCPInfo`
    FOREIGN KEY (`Test_TCPInfo_id`)
    REFERENCES `monitorDB`.`TCPInfo` (`TCPInfo_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_Tests_Service`
    FOREIGN KEY (`Test_Service_id`)
    REFERENCES `monitorDB`.`Service` (`Service_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION,
  CONSTRAINT `fk_Tests_DaySlice`
    FOREIGN KEY (`Test_DaySlice_id`)
    REFERENCES `monitorDB`.`DaySlice` (`DaySlice_id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION)
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `monitorDB`.`DayStat`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `monitorDB`.`DayStat` (
  `DayStat_id` INT NOT NULL AUTO_INCREMENT,
  `DayStat_Type` VARCHAR(45) NULL,
  `DayStat_Date` DATE NULL,
  `DayStat_AvgBW` INT NULL,
  `DayStat_MinBW` INT NULL,
  `DayStat_MaxBW` INT NULL,
  `DayStat_MedianBW` INT NULL,
  `DayStat_AvgMinRTT` INT NULL,
  `DayStat_MinMinRTT` INT NULL,
  `DayStat_MaxMinRTT` INT NULL,
  `DayStat_MedianMinRTT` INT NULL,
  `DayStat_Down_Number` INT NULL,
  `DayStat_Up_Number` INT NULL,
  `DayStat_Total` INT NULL,
  PRIMARY KEY (`DayStat_id`))
ENGINE = InnoDB;


-- -----------------------------------------------------
-- Table `monitorDB`.`InterDayStat`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `monitorDB`.`InterDayStat` (
  `Inter_DayStat_id` INT NOT NULL AUTO_INCREMENT,
  `Inter_DayStat_Type` VARCHAR(45) NULL,
  `Inter_DayStat_Date` DATE NULL,
  `Inter_DayStat_AvgBW` INT NULL,
  `Inter_DayStat_MinBW` INT NULL,
  `Inter_DayStat_MaxBW` INT NULL,
  `Inter_DayStat_MedianBW` INT NULL,
  `Inter_DayStat_AvgMinRTT` INT NULL,
  `Inter_DayStat_MinMinRTT` INT NULL,
  `Inter_DayStat_MaxMinRTT` INT NULL,
  `Inter_DayStat_MedianMinRTT` INT NULL,
  `Inter_DayStat_Down_Number` INT NULL,
  `Inter_DayStat_Up_Number` INT NULL,
  `Inter_DayStat_Total` INT NULL,
  PRIMARY KEY (`Inter_DayStat_id`))
ENGINE = InnoDB;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
