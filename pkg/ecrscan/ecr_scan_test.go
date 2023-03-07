package ecrscan

import (
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/aws/aws-sdk-go/service/ecr/ecriface"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

type mockECRClient struct {
	mock.Mock
	ecriface.ECRAPI
}

const maxScanAge = 24

var logger, _ = zap.NewProduction()
var ecrClient = &mockECRClient{}
var evaluator = Evaluator{
	MaxScanAge: maxScanAge,
	Logger:     logger,
	ECRClient:  ecrClient,
}

var testCases = map[string]*ecr.DescribeImageScanFindingsOutput{
	"ScanCompletedNoFindings": {
		ImageScanFindings: &ecr.ImageScanFindings{
			FindingSeverityCounts: map[string]*int64{
				ecr.FindingSeverityUndefined:     aws.Int64(0),
				ecr.FindingSeverityInformational: aws.Int64(0),
				ecr.FindingSeverityLow:           aws.Int64(0),
				ecr.FindingSeverityMedium:        aws.Int64(0),
				ecr.FindingSeverityHigh:          aws.Int64(0),
				ecr.FindingSeverityCritical:      aws.Int64(0),
			},
			ImageScanCompletedAt: relativeTimePointer(1),
		},
		ImageScanStatus: &ecr.ImageScanStatus{
			Status: aws.String(ecr.ScanStatusComplete),
		},
	},
	"ScanCompletedOneUndefinedFinding": {
		ImageScanFindings: &ecr.ImageScanFindings{
			FindingSeverityCounts: map[string]*int64{
				ecr.FindingSeverityUndefined:     aws.Int64(1),
				ecr.FindingSeverityInformational: aws.Int64(0),
				ecr.FindingSeverityLow:           aws.Int64(0),
				ecr.FindingSeverityMedium:        aws.Int64(0),
				ecr.FindingSeverityHigh:          aws.Int64(0),
				ecr.FindingSeverityCritical:      aws.Int64(0),
			},
			ImageScanCompletedAt: relativeTimePointer(1),
		},
		ImageScanStatus: &ecr.ImageScanStatus{
			Status: aws.String(ecr.ScanStatusComplete),
		},
	},
	"ScanCompletedOneCriticalFinding": {
		ImageScanFindings: &ecr.ImageScanFindings{
			FindingSeverityCounts: map[string]*int64{
				ecr.FindingSeverityUndefined:     aws.Int64(0),
				ecr.FindingSeverityInformational: aws.Int64(0),
				ecr.FindingSeverityLow:           aws.Int64(0),
				ecr.FindingSeverityMedium:        aws.Int64(0),
				ecr.FindingSeverityHigh:          aws.Int64(0),
				ecr.FindingSeverityCritical:      aws.Int64(1),
			},
			ImageScanCompletedAt: relativeTimePointer(1),
		},
		ImageScanStatus: &ecr.ImageScanStatus{
			Status: aws.String(ecr.ScanStatusComplete),
		},
	},

	"ScanCompletedOneFindingEachCategory": {
		ImageScanFindings: &ecr.ImageScanFindings{
			FindingSeverityCounts: map[string]*int64{
				ecr.FindingSeverityUndefined:     aws.Int64(1),
				ecr.FindingSeverityInformational: aws.Int64(1),
				ecr.FindingSeverityLow:           aws.Int64(1),
				ecr.FindingSeverityMedium:        aws.Int64(1),
				ecr.FindingSeverityHigh:          aws.Int64(1),
				ecr.FindingSeverityCritical:      aws.Int64(1),
			},
			ImageScanCompletedAt: relativeTimePointer(1),
		},
		ImageScanStatus: &ecr.ImageScanStatus{
			Status: aws.String(ecr.ScanStatusComplete),
		},
	},
	"ScanCompletedMultipleFindingsEachCategory": {
		ImageScanFindings: &ecr.ImageScanFindings{
			FindingSeverityCounts: map[string]*int64{
				ecr.FindingSeverityUndefined:     aws.Int64(5),
				ecr.FindingSeverityInformational: aws.Int64(8),
				ecr.FindingSeverityLow:           aws.Int64(13),
				ecr.FindingSeverityMedium:        aws.Int64(21),
				ecr.FindingSeverityHigh:          aws.Int64(34),
				ecr.FindingSeverityCritical:      aws.Int64(55),
			},
			ImageScanCompletedAt: relativeTimePointer(1),
		},
		ImageScanStatus: &ecr.ImageScanStatus{
			Status: aws.String(ecr.ScanStatusComplete),
		},
	},
}

// relativeTimePointer is a helper function that returns a pointer to a time
// object relative to the current time. A negative number represents a time in
// the future; a positive number represents a time in the past.
//
// Example:
//
//	relativeTimePointer(-3) -- 3 hours from now
//	relativeTimePointer(1)  -- 1 hour ago
//	relativeTimePointer(0)  -- now
func relativeTimePointer(hours float64) *time.Time {
	t := time.Now()
	newT := t.Add(-time.Duration(hours) * time.Hour)
	return &newT
}

func (m *mockECRClient) DescribeImageScanFindings(input *ecr.DescribeImageScanFindingsInput) (*ecr.DescribeImageScanFindingsOutput, error) {
	if _, ok := testCases[*input.ImageId.ImageTag]; ok {
		return testCases[*input.ImageId.ImageTag], nil
	}
	return nil, errors.New("error")
}

func (m *mockECRClient) WaitUntilImageScanComplete(input *ecr.DescribeImageScanFindingsInput) error {
	return nil
}

func TestEvaluate(t *testing.T) {
	tests := []struct {
		target   *Target
		expected Report
	}{
		{
			&Target{
				ImageTag:   "ScanCompletedNoFindings",
				Repository: "test-repo",
			},
			Report{
				TotalFindings: 0,
			},
		},
		{
			&Target{
				ImageTag:   "ScanCompletedOneUndefinedFinding",
				Repository: "test-repo",
			},
			Report{
				TotalFindings: 1,
			},
		},
		{
			&Target{
				ImageTag:   "ScanCompletedOneCriticalFinding",
				Repository: "test-repo",
			},
			Report{
				TotalFindings: 1,
			},
		},
		{
			&Target{
				ImageTag:   "ScanCompletedOneFindingEachCategory",
				Repository: "test-repo",
			},
			Report{
				TotalFindings: 6,
			},
		},
		{
			&Target{
				ImageTag:   "ScanCompletedMultipleFindingsEachCategory",
				Repository: "test-repo",
			},
			Report{
				TotalFindings: 136,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.target.ImageTag, func(t *testing.T) {
			report, err := evaluator.Evaluate(tt.target)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, report)
		})
	}
}

func TestEvaluateWithBadInput(t *testing.T) {
	tests := []struct {
		description string
		target      *Target
	}{
		{
			"Nil target",
			nil,
		},
		{
			"Empty target",
			&Target{},
		},
		{
			"No repository",
			&Target{
				ImageTag: "test123",
			},
		},
		{
			"No image tag",
			&Target{
				Repository: "testrepo",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			_, err := evaluator.Evaluate(tt.target)
			assert.NotNil(t, err)
		})
	}
}

func TestIsOldScan(t *testing.T) {
	tests := []struct {
		description string
		findings    *ecr.ImageScanFindings
		expected    bool
	}{
		{
			"Scan created in the future",
			&ecr.ImageScanFindings{
				ImageScanCompletedAt: relativeTimePointer(-3),
			},
			false,
		},
		{
			"Scan created now",
			&ecr.ImageScanFindings{
				ImageScanCompletedAt: relativeTimePointer(0),
			},
			false,
		},
		{
			"Scan created 1 hour ago",
			&ecr.ImageScanFindings{
				ImageScanCompletedAt: relativeTimePointer(1),
			},
			false,
		},
		{
			"Scan created 1 hour before max age",
			&ecr.ImageScanFindings{
				ImageScanCompletedAt: relativeTimePointer(float64(maxScanAge - 1)),
			},
			false,
		},
		{
			"Scan created 1 hour after max age",
			&ecr.ImageScanFindings{
				ImageScanCompletedAt: relativeTimePointer(float64(maxScanAge + 1)),
			},
			true,
		},
		{
			"Scan created fraction before max age",
			&ecr.ImageScanFindings{
				ImageScanCompletedAt: relativeTimePointer(float64(maxScanAge) - 0.0000000000001),
			},
			false,
		},
		{
			"Scan created fraction after max age",
			&ecr.ImageScanFindings{
				ImageScanCompletedAt: relativeTimePointer(float64(maxScanAge) + 0.0000000000001),
			},
			true,
		},
		{
			"Scan created far in the past",
			&ecr.ImageScanFindings{
				ImageScanCompletedAt: relativeTimePointer(float64(maxScanAge + 999999)),
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			assert.Equal(t, tt.expected, evaluator.isOldScan(tt.findings.ImageScanCompletedAt))
		})
	}
}

func TestCalculateTotalFindingsWithBadInput(t *testing.T) {
	tests := []struct {
		description string
		findings    *ecr.ImageScanFindings
		expected    int
	}{
		{
			"Nil scan findings",
			nil,
			-1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			result, err := evaluator.calculateTotalFindings(tt.findings)
			assert.Equal(t, tt.expected, result, "result should be -1")
			assert.NotNil(t, err)
		})
	}

}
