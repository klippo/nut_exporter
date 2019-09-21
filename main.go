package main

import (
	"strconv"
	"os/exec"
 	"time"
	"fmt"
	"flag"
	"regexp"
	"net/http"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
)
// Regex
var(
	batteryChargeRegex          =   regexp.MustCompile(`(?:battery[.]charge:(?:\s)(.*))`)
	batteryChargeLowRegex       =   regexp.MustCompile(`(?:battery[.]charge.low:(?:\s)(.*))`)
	batteryChargeWarningRegex   =   regexp.MustCompile(`(?:battery[.]charge.warning:(?:\s)(.*))`)
	batteryPacksRegex           =   regexp.MustCompile(`(?:battery[.]packs:(?:\s)(.*))`)
	batteryRuntimeRegex         =   regexp.MustCompile(`(?:battery[.]runtime:(?:\s)(.*))`)
	batteryRuntimeLowRegex      =   regexp.MustCompile(`(?:battery[.]runtime.low:(?:\s)(.*))`)
	batteryTemperatureRegex     =   regexp.MustCompile(`(?:battery[.]temperature:(?:\s)(.*))`)
	batteryVoltageRegex         =   regexp.MustCompile(`(?:battery[.]voltage:(?:\s)(.*))`)
	batteryVoltageNominalRegex  =   regexp.MustCompile(`(?:battery[.]voltage[.]nominal:(?:\s)(.*))`)
	inputTransferLowRegex       =   regexp.MustCompile(`(?:input[.]transfer.low:(?:\s)(.*))`)
	inputTransferHighRegex      =   regexp.MustCompile(`(?:input[.]transfer.high:(?:\s)(.*))`)
	inputVoltageRegex           =   regexp.MustCompile(`(?:input[.]voltage:(?:\s)(.*))`)
	inputVoltageNominalRegex    =   regexp.MustCompile(`(?:input[.]voltage[.]nominal:(?:\s)(.*))`)
	outputCurrentRegex          =   regexp.MustCompile(`(?:output[.]current:(?:\s)(.*))`)
	outputFrequencyRegex        =   regexp.MustCompile(`(?:output[.]frequency:(?:\s)(.*))`)
	outputVoltageRegex          =   regexp.MustCompile(`(?:output[.]voltage:(?:\s)(.*))`)
	outputVoltageNominalRegex   =   regexp.MustCompile(`(?:output[.]voltage[.]nominal:(?:\s)(.*))`)
	upsLoadRegex                =   regexp.MustCompile(`(?:ups[.]load:(?:\s)(.*))`)
	upsPowerNominalRegex        =   regexp.MustCompile(`(?:ups[.]power[.]nominal:(?:\s)(.*))`)
	upsRealPowerNominalRegex        =   regexp.MustCompile(`(?:ups[.]realpower[.]nominal:(?:\s)(.*))`)
	upsStatusRegex              =   regexp.MustCompile(`(?:ups[.]status:(?:\s)(.*))`)
	upsTempRegex                =   regexp.MustCompile(`(?:ups[.]temperature:(?:\s)(.*))`)
)
// NUT Gauges
var (
	batteryCharge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_battery_charge",
		Help: "Current battery charge (percent)",
	})

	batteryChargeLow = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_battery_charge_low",
		Help: "Current battery charge (percent)",
	})

	batteryChargeWarning = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_battery_charge_warning",
		Help: "Current battery charge (percent)",
	})

	batteryPacks = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_battery_pack",
		Help: "Number of battery packs on the UPS",
	})

	batteryRuntime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_battery_runtime",
		Help: "Current battery charge (percent)",
	})

	batteryRuntimeLow = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_battery_runtime_low",
		Help: "Current battery charge (percent)",
	})

	batteryTemperature = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_battery_temperature",
		Help: "Current battery charge (percent)",
	})

	batteryVoltage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_battery_voltage",
		Help: "Current battery voltage",
	})

	batteryVoltageNominal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_battery_voltage_nominal",
		Help: "Nominal battery voltage",
	})

	inputTransferLow = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_input_transfer_low",
		Help: "Current input voltage",
	})
	inputTransferHigh = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_input_transfer_high",
		Help: "Current input voltage",
	})

	inputVoltage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_input_voltage",
		Help: "Current input voltage",
	})

	inputVoltageNominal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_input_voltage_nominal",
		Help: "Nominal input voltage",
	})

	outputCurrent = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_output_current",
		Help: "Current output voltage",
	})
	outputFrequency = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_output_frequency",
		Help: "Current output voltage",
	})
	outputVoltage = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_output_voltage",
		Help: "Current output voltage",
	})
	
	outputVoltageNominal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_output_voltage_nominal",
		Help: "Nominal output voltage",
	})
	
	upsPowerNominal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_ups_power_nominal",
		Help: "Nominal ups power",
	})
	
	upsRealPowerNominal = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_ups_realpower_nominal",
		Help: "Nominal ups realpower",
	})
	
	upsTemp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_ups_temp",
		Help: "UPS Temperature (degrees C)",
	})
	
	upsLoad = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_ups_load",
		Help: "Current UPS load (percent)",
	})

	upsStatus = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "nut_ups_status",
		Help: "Current UPS Status (0=Calibration, 1=SmartTrim, 2=SmartBoost, 3=Online, 4=OnBattery, 5=Overloaded, 6=LowBattery, 7=ReplaceBattery, 8=OnBypass, 9=Off, 10=Charging, 11=Discharging)",
	})
)

func recordMetrics(upscBinary string, upsArg string){
	prometheus.MustRegister(batteryCharge)
	prometheus.MustRegister(batteryChargeLow)
	prometheus.MustRegister(batteryChargeWarning)
	prometheus.MustRegister(batteryPacks)
	prometheus.MustRegister(batteryRuntime)
	prometheus.MustRegister(batteryRuntimeLow)
	prometheus.MustRegister(batteryTemperature)
	prometheus.MustRegister(batteryVoltage)
	prometheus.MustRegister(batteryVoltageNominal)
	prometheus.MustRegister(inputTransferLow)
	prometheus.MustRegister(inputTransferHigh)
	prometheus.MustRegister(inputVoltage)
	prometheus.MustRegister(inputVoltageNominal)
	prometheus.MustRegister(outputCurrent)
	prometheus.MustRegister(outputFrequency)
	prometheus.MustRegister(outputVoltage)
	prometheus.MustRegister(outputVoltageNominal)
	prometheus.MustRegister(upsLoad)
	prometheus.MustRegister(upsPowerNominal)
	prometheus.MustRegister(upsRealPowerNominal)
	prometheus.MustRegister(upsStatus)
	prometheus.MustRegister(upsTemp)

	go func(){
		for {
			upsOutput, err := exec.Command(upscBinary , upsArg).Output()
			
			if err != nil {
				log.Fatal(err)
			}
			
			if batteryChargeRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(batteryCharge)
			} else {
				batteryChargeValue, _ := strconv.ParseFloat(batteryChargeRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				batteryCharge.Set(batteryChargeValue)
			}

			if batteryChargeLowRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(batteryChargeLow)
			} else {
				batteryChargeLowValue, _ := strconv.ParseFloat(batteryChargeLowRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				batteryChargeLow.Set(batteryChargeLowValue)

			}
			if batteryChargeWarningRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(batteryChargeWarning)
			} else {
				batteryChargeWarningValue, _ := strconv.ParseFloat(batteryChargeWarningRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				batteryChargeWarning.Set(batteryChargeWarningValue)
			}

			if batteryPacksRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(batteryPacks)
			} else {
				batteryPacksValue, _ := strconv.ParseFloat(batteryPacksRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				batteryPacks.Set(batteryPacksValue)
			}

			if batteryRuntimeRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(batteryRuntime)
			} else {
				batteryRuntimeValue, _ := strconv.ParseFloat(batteryRuntimeRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				batteryRuntime.Set(batteryRuntimeValue)
			}

			if batteryRuntimeLowRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(batteryRuntimeLow)
			} else {
				batteryRuntimeLowValue, _ := strconv.ParseFloat(batteryRuntimeLowRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				batteryRuntimeLow.Set(batteryRuntimeLowValue)
			}

			if batteryTemperatureRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(batteryTemperature)
			} else {
				batteryTemperatureValue, _ := strconv.ParseFloat(batteryTemperatureRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				batteryTemperature.Set(batteryTemperatureValue)
			}

			if batteryVoltageRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(batteryVoltage)
			} else {
				batteryVoltageValue, _ := strconv.ParseFloat(batteryVoltageRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				batteryVoltage.Set(batteryVoltageValue)
			}
			
			if batteryVoltageNominalRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(batteryVoltageNominal)
			} else {
				batteryVoltageNominalValue, _ := strconv.ParseFloat(batteryVoltageNominalRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				batteryVoltageNominal.Set(batteryVoltageNominalValue)
			}

			if inputTransferLowRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(inputTransferLow)
			} else {
				inputTransferLowValue, _ := strconv.ParseFloat(inputTransferLowRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				inputTransferLow.Set(inputTransferLowValue)
			}

			if inputTransferHighRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(inputTransferHigh)
			} else {
				inputTransferHighValue, _ := strconv.ParseFloat(inputTransferHighRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				inputTransferHigh.Set(inputTransferHighValue)
			}

			if inputVoltageRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(inputVoltage)
			} else {
				inputVoltageValue, _ := strconv.ParseFloat(inputVoltageRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				inputVoltage.Set(inputVoltageValue)
			}
			
			if inputVoltageNominalRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(inputVoltageNominal)
			} else {
				inputVoltageNominalValue, _ := strconv.ParseFloat(inputVoltageNominalRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				inputVoltageNominal.Set(inputVoltageNominalValue)
			}

			if outputCurrentRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(outputCurrent)
			} else {
				outputCurrentValue, _ := strconv.ParseFloat(outputCurrentRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				outputCurrent.Set(outputCurrentValue)
			}

			if outputFrequencyRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(outputFrequency)
			} else {
				outputFrequencyValue, _ := strconv.ParseFloat(outputFrequencyRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				outputFrequency.Set(outputFrequencyValue)
			}

			if outputVoltageRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(outputVoltage)
			} else {
				outputVoltageValue, _ := strconv.ParseFloat(outputVoltageRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
  				outputVoltage.Set(outputVoltageValue)
			}
			
			if outputVoltageNominalRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(outputVoltageNominal)
			} else {
				outputVoltageNominalValue, _ := strconv.ParseFloat(outputVoltageNominalRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				outputVoltageNominal.Set(outputVoltageNominalValue)
			}
			
			if upsPowerNominalRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(upsPowerNominal)
			} else {
				upsPowerNominalValue, _ := strconv.ParseFloat(upsPowerNominalRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				upsPowerNominal.Set(upsPowerNominalValue)
			}
			
			if upsRealPowerNominalRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(upsRealPowerNominal)
			} else {
				upsRealPowerNominalValue, _ := strconv.ParseFloat(upsRealPowerNominalRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)
				upsRealPowerNominal.Set(upsRealPowerNominalValue)
			}
			
			if upsTempRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(upsTemp)
			} else {
				upsTempValue, _ := strconv.ParseFloat(upsTempRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)	
				upsTemp.Set(upsTempValue)	
			}
			
			if upsLoadRegex.FindAllStringSubmatch(string(upsOutput), -1) == nil {
				prometheus.Unregister(upsLoad)
			} else {
				upsLoadValue, _ := strconv.ParseFloat(upsLoadRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1], 64)	
				upsLoad.Set(upsLoadValue)
			}
			
			upsStatusValue := upsStatusRegex.FindAllStringSubmatch(string(upsOutput), -1)[0][1]	
			
			switch upsStatusValue {
				case "CAL":
            		upsStatus.Set(0)
				case "TRIM":
					upsStatus.Set(1)
				case "BOOST":
					upsStatus.Set(2)
				case "OL":
					upsStatus.Set(3)
				case "OB":
					upsStatus.Set(4)
				case "OVER":
					upsStatus.Set(5)
				case "LB":
					upsStatus.Set(6)
				case "RB":
					upsStatus.Set(7)
				case "BYPASS":
					upsStatus.Set(8)
				case "OFF":
					upsStatus.Set(9)
				case "CHRG":
					upsStatus.Set(10)
				case "DISCHRG":
					upsStatus.Set(11)
			}
			time.Sleep(5 * time.Second)
		}
	}()
}

func main() {
	upsArg   := flag.String("ups", "none", "ups name managed by nut")
	portArg  := flag.Int("port", 8100, "port number")
	upscArg  := flag.String("upsc", "/bin/upsc", "upsc path")
    flag.Parse()

	var listenAddr = fmt.Sprintf(":%d", *portArg)
	recordMetrics(*upscArg, *upsArg)
    
	log.Infoln("Starting NUT exporter on ups", *upsArg )	
	http.Handle("/metrics", promhttp.Handler())
    
	log.Infoln("NUT exporter started on port", *portArg)
	http.ListenAndServe(listenAddr, nil)	
}
