package seed

import (
	"math/rand"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/yourname/payslip-system/internal/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Run(db *gorm.DB) error {
	rand.Seed(time.Now().UnixNano())

	//create 100 employees
	for i := 0; i < 100; i++ {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		employee := model.Employee{
			Name:     faker.Name(),
			Password: string(hashedPassword),
			Role:     "employee",
			Active:   true,
		}
		if err := db.Create(&employee).Error; err != nil {
			return err
		}
	}

	//create 1 admin
	adminPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := model.Employee{
		Name:     "Admin",
		Password: string(adminPassword),
		Role:     "admin",
		Active:   true,
	}
	if err := db.Create(&admin).Error; err != nil {
		return err
	}

	// Generate random attendance for one month (current month)
	// if err := generateRandomAttendance(db); err != nil {
	// 	return err
	// }

	// Generate random overtime records
	// if err := generateRandomOvertimes(db); err != nil {
	// 	return err
	// }

	// Generate random reimbursement records
	// if err := generateRandomReimbursements(db); err != nil {
	// 	return err
	// }

	return nil
}

func generateRandomAttendance(db *gorm.DB) error {
	// Get all employees
	var employees []model.Employee
	if err := db.Find(&employees).Error; err != nil {
		return err
	}

	// Get current month start and end dates
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	// Generate attendance for each employee for each working day in the month
	for _, employee := range employees {
		for date := startOfMonth; !date.After(endOfMonth); date = date.AddDate(0, 0, 1) {
			// Skip weekends
			if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
				continue
			}

			// 90% chance of attendance (some employees might be absent)
			if rand.Float32() > 0.9 {
				// Create absent record
				attendance := model.Attendance{
					EmployeeID:  employee.ID,
					Status:      "absent",
					Date:        date,
					HoursWorked: 0,
				}
				if err := db.Create(&attendance).Error; err != nil {
					return err
				}
				continue
			}

			// Generate random checkin time (8:00 AM to 9:30 AM)
			checkinHour := 8 + rand.Intn(2) // 8 or 9
			checkinMinute := rand.Intn(60)  // 0-59 minutes
			if checkinHour == 9 && checkinMinute > 30 {
				checkinMinute = 30 // Cap at 9:30
			}

			checkin := time.Date(date.Year(), date.Month(), date.Day(),
				checkinHour, checkinMinute, 0, 0, date.Location())

			// Generate random checkout time (4:30 PM to 6:00 PM)
			checkoutHour := 16 + rand.Intn(3) // 16, 17, or 18
			checkoutMinute := rand.Intn(60)   // 0-59 minutes
			if checkoutHour == 16 && checkoutMinute < 30 {
				checkoutMinute = 30 // Minimum 4:30 PM
			}
			if checkoutHour == 18 && checkoutMinute > 0 {
				checkoutHour = 17
				checkoutMinute = 59 // Cap at 5:59 PM
			}

			checkout := time.Date(date.Year(), date.Month(), date.Day(),
				checkoutHour, checkoutMinute, 0, 0, date.Location())

			// 5% chance of half-day (early checkout)
			if rand.Float32() < 0.05 {
				checkout = time.Date(date.Year(), date.Month(), date.Day(),
					13, rand.Intn(60), 0, 0, date.Location()) // 1:00-1:59 PM
			}

			// 3% chance of overtime (late checkout)
			if rand.Float32() < 0.03 {
				overtimeHour := 19 + rand.Intn(3) // 7-9 PM
				checkout = time.Date(date.Year(), date.Month(), date.Day(),
					overtimeHour, rand.Intn(60), 0, 0, date.Location())
			}

			// Create attendance record
			attendance := model.Attendance{
				EmployeeID: employee.ID,
				Checkin:    checkin,
				Checkout:   &checkout,
				Status:     "present",
				Date:       date,
			}

			// Calculate hours worked
			attendance.CalculateHours()

			if err := db.Create(&attendance).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func generateRandomOvertimes(db *gorm.DB) error {
	// Get all employees (excluding admin for overtime requests)
	var employees []model.Employee
	if err := db.Where("role = ?", "employee").Find(&employees).Error; err != nil {
		return err
	}

	// Get admin for approval
	var admin model.Employee
	if err := db.Where("role = ?", "admin").First(&admin).Error; err != nil {
		return err
	}

	// Get current month start and end dates
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	// Predefined overtime reasons
	overtimeReasons := []string{
		"Project deadline requirements",
		"Client meeting preparation",
		"System maintenance work",
		"Urgent bug fixes and testing",
		"Monthly report compilation",
		"Database backup and optimization",
		"Training session preparation",
		"Critical client support",
		"End-of-month financial closing",
		"Software deployment activities",
		"Emergency incident response",
		"Presentation preparation for board meeting",
	}

	// Generate overtime records for employees
	for _, employee := range employees {
		// Each employee has 10-30% chance of overtime days in a month (including weekends/holidays)
		totalDaysInMonth := endOfMonth.Day()

		// Determine number of overtime days (10-30% of total days in month)
		minOvertimeDays := int(float64(totalDaysInMonth) * 0.1)
		maxOvertimeDays := int(float64(totalDaysInMonth) * 0.3)
		if minOvertimeDays == 0 {
			minOvertimeDays = 1
		}
		if maxOvertimeDays <= minOvertimeDays {
			maxOvertimeDays = minOvertimeDays + 1
		}

		overtimeDays := minOvertimeDays + rand.Intn(maxOvertimeDays-minOvertimeDays+1)

		// Generate random overtime dates (can be workdays or holidays/weekends)
		var selectedDates []time.Time
		for len(selectedDates) < overtimeDays {
			randomDay := 1 + rand.Intn(endOfMonth.Day())
			randomDate := time.Date(now.Year(), now.Month(), randomDay, 0, 0, 0, 0, now.Location())

			// Check if date already selected
			dateExists := false
			for _, existing := range selectedDates {
				if existing.Equal(randomDate) {
					dateExists = true
					break
				}
			}

			if !dateExists && !randomDate.After(now) {
				selectedDates = append(selectedDates, randomDate)
			}
		}

		// Create overtime records for selected dates
		for _, overtimeDate := range selectedDates {
			// Maximum 3 hours overtime per day
			hours := 1 + rand.Intn(3) // 1-3 hours

			var startTime time.Time
			var endTime time.Time

			// Check if it's a workday or holiday/weekend
			isWorkday := overtimeDate.Weekday() != time.Saturday && overtimeDate.Weekday() != time.Sunday

			if isWorkday {
				// On workdays: overtime starts after checkout (after 5:00 PM minimum)
				// Get employee's checkout time for this date or use default 5:00 PM
				var attendance model.Attendance
				checkoutTime := time.Date(overtimeDate.Year(), overtimeDate.Month(), overtimeDate.Day(), 17, 0, 0, 0, overtimeDate.Location()) // Default 5:00 PM

				// Try to get actual checkout time from attendance
				if err := db.Where("employee_id = ? AND date = ?", employee.ID, overtimeDate).First(&attendance).Error; err == nil && attendance.Checkout != nil {
					checkoutTime = *attendance.Checkout
				}

				// Overtime starts 15-60 minutes after checkout
				minutesAfterCheckout := 15 + rand.Intn(46) // 15-60 minutes
				startTime = checkoutTime.Add(time.Duration(minutesAfterCheckout) * time.Minute)
			} else {
				// On holidays/weekends: overtime can start anytime (9 AM to 6 PM)
				startHour := 9 + rand.Intn(10) // 9 AM to 6 PM
				startMinute := rand.Intn(60)
				startTime = time.Date(overtimeDate.Year(), overtimeDate.Month(), overtimeDate.Day(),
					startHour, startMinute, 0, 0, overtimeDate.Location())
			}

			// Calculate end time (start time + hours)
			endTime = startTime.Add(time.Duration(hours) * time.Hour)

			// Random reason
			reason := overtimeReasons[rand.Intn(len(overtimeReasons))]

			// Random status distribution: 70% approved, 20% pending, 10% rejected
			var status model.OvertimeStatus
			var approvedBy *uint
			var approvedAt *time.Time

			statusRand := rand.Float32()
			if statusRand < 0.7 {
				// Approved
				status = model.OvertimeApproved
				approvedBy = &admin.ID
				approvalTime := overtimeDate.AddDate(0, 0, rand.Intn(3)+1) // Approved 1-3 days later
				approvedAt = &approvalTime
			} else if statusRand < 0.9 {
				// Pending
				status = model.OvertimePending
			} else {
				// Rejected
				status = model.OvertimeRejected
				approvedBy = &admin.ID
				rejectionTime := overtimeDate.AddDate(0, 0, rand.Intn(3)+1) // Rejected 1-3 days later
				approvedAt = &rejectionTime
			}

			// Create overtime record
			overtime := model.Overtime{
				EmployeeID:   employee.ID,
				OvertimeDate: overtimeDate.Format("2006-01-02"), // Convert time.Time to string in "YYYY-MM-DD" format
				StartTime:    startTime.Format("15:04:05"),
				EndTime:      endTime.Format("15:04:05"),
				Hours:        hours,
				Reason:       reason,
				Status:       status,
				ApprovedBy:   approvedBy,
				ApprovedAt:   approvedAt,
			}
			// If OvertimeDate is string in your model, use:
			// OvertimeDate: overtimeDate.Format("2006-01-02"),

			if err := db.Create(&overtime).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func generateRandomReimbursements(db *gorm.DB) error {
	// Get all employees (excluding admin for reimbursement requests)
	var employees []model.Employee
	if err := db.Where("role = ?", "employee").Find(&employees).Error; err != nil {
		return err
	}

	// Get admin for approval
	var admin model.Employee
	if err := db.Where("role = ?", "admin").First(&admin).Error; err != nil {
		return err
	}

	// Get current month start and end dates
	now := time.Now()
	startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, -1)

	// Predefined reimbursement data by category
	reimbursementData := map[model.ReimbursementCategory]struct {
		reasons     []string
		minAmount   float64
		maxAmount   float64
		probability float32 // probability of this category being selected
	}{
		model.ReimbursementTravel: {
			reasons: []string{
				"Business trip to client office",
				"Conference attendance transportation",
				"Project site visit expenses",
				"Client meeting travel costs",
				"Training seminar travel",
				"Airport transfer and parking fees",
			},
			minAmount:   15.00,
			maxAmount:   500.00,
			probability: 0.25,
		},
		model.ReimbursementMeals: {
			reasons: []string{
				"Client lunch meeting",
				"Business dinner with partners",
				"Conference meal expenses",
				"Overtime meal allowance",
				"Team building dinner",
				"Working lunch during project",
			},
			minAmount:   10.00,
			maxAmount:   150.00,
			probability: 0.35,
		},
		model.ReimbursementEquipment: {
			reasons: []string{
				"Office supplies for remote work",
				"Software license for project",
				"Ergonomic keyboard and mouse",
				"Monitor for better productivity",
				"Webcam for video conferences",
				"Headset for client calls",
			},
			minAmount:   25.00,
			maxAmount:   800.00,
			probability: 0.15,
		},
		model.ReimbursementTraining: {
			reasons: []string{
				"Professional certification course",
				"Technical training workshop",
				"Industry conference registration",
				"Online course subscription",
				"Skills development program",
				"Professional seminar attendance",
			},
			minAmount:   50.00,
			maxAmount:   1200.00,
			probability: 0.10,
		},
		model.ReimbursementMedical: {
			reasons: []string{
				"Annual health checkup",
				"Prescription glasses for work",
				"Medical consultation fees",
				"Vaccination for business travel",
				"Work-related injury treatment",
				"Eye examination for computer work",
			},
			minAmount:   20.00,
			maxAmount:   600.00,
			probability: 0.10,
		},
		model.ReimbursementOther: {
			reasons: []string{
				"Internet upgrade for remote work",
				"Phone bill for business calls",
				"Home office setup costs",
				"Parking fees for client visits",
				"Document printing and binding",
				"Business card printing",
			},
			minAmount:   5.00,
			maxAmount:   200.00,
			probability: 0.05,
		},
	}
	// Generate reimbursement records for employees
	for _, employee := range employees {
		// Each employee has 20-60% chance of reimbursement requests in a month
		// Determine number of reimbursement requests (1-8 per month)
		maxReimbursements := 1 + rand.Intn(8) // 1-8 requests per month
		reimbursementCount := 1 + rand.Intn(maxReimbursements)

		// Generate random reimbursement dates
		var selectedDates []time.Time
		for len(selectedDates) < reimbursementCount {
			randomDay := 1 + rand.Intn(endOfMonth.Day())
			randomDate := time.Date(now.Year(), now.Month(), randomDay, 0, 0, 0, 0, now.Location())

			// Check if date already selected
			dateExists := false
			for _, existing := range selectedDates {
				if existing.Equal(randomDate) {
					dateExists = true
					break
				}
			}

			if !dateExists && !randomDate.After(now) {
				selectedDates = append(selectedDates, randomDate)
			}
		}

		// Create reimbursement records for selected dates
		for _, reimbursementDate := range selectedDates {
			// Select category based on probability
			categoryRand := rand.Float32()
			var selectedCategory model.ReimbursementCategory
			var categoryData struct {
				reasons     []string
				minAmount   float64
				maxAmount   float64
				probability float32
			}

			cumulative := float32(0.0)
			for category, data := range reimbursementData {
				cumulative += data.probability
				if categoryRand <= cumulative {
					selectedCategory = category
					categoryData = data
					break
				}
			}

			// If no category selected (shouldn't happen), default to meals
			if selectedCategory == "" {
				selectedCategory = model.ReimbursementMeals
				categoryData = reimbursementData[model.ReimbursementMeals]
			}

			// Generate random amount within category range
			amountRange := categoryData.maxAmount - categoryData.minAmount
			amount := categoryData.minAmount + (rand.Float64() * amountRange)
			// Round to 2 decimal places
			amount = float64(int(amount*100)) / 100

			// Select random reason from category
			reason := categoryData.reasons[rand.Intn(len(categoryData.reasons))]

			// Random status distribution: 60% approved, 15% paid, 20% pending, 5% rejected
			var status model.ReimbursementStatus
			var approvedBy *uint
			var approvedAt *time.Time

			statusRand := rand.Float32()
			if statusRand < 0.6 {
				// Approved
				status = model.ReimbursementApproved
				approvedBy = &admin.ID
				approvalTime := reimbursementDate.AddDate(0, 0, rand.Intn(5)+1) // Approved 1-5 days later
				approvedAt = &approvalTime
			} else if statusRand < 0.75 {
				// Paid (automatically approved first)
				status = model.ReimbursementPaid
				approvedBy = &admin.ID
				approvalTime := reimbursementDate.AddDate(0, 0, rand.Intn(3)+1) // Approved 1-3 days later
				approvedAt = &approvalTime
			} else if statusRand < 0.95 {
				// Pending
				status = model.ReimbursementPending
			} else {
				// Rejected
				status = model.ReimbursementRejected
				approvedBy = &admin.ID
				rejectionTime := reimbursementDate.AddDate(0, 0, rand.Intn(7)+1) // Rejected 1-7 days later
				approvedAt = &rejectionTime
			}

			// Create reimbursement record
			reimbursement := model.Reimbursement{
				EmployeeID:        employee.ID,
				ReimbursementDate: reimbursementDate,
				Amount:            amount,
				Category:          selectedCategory,
				Reason:            reason,
				Status:            status,
				ApprovedBy:        approvedBy,
				ApprovedAt:        approvedAt,
			}

			if err := db.Create(&reimbursement).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
