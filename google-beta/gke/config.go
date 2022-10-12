package gke

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"time"

	compute "google.golang.org/api/compute/v0.beta"
	container "google.golang.org/api/container/v1beta1"
	"google.golang.org/api/option"
)

type providerMeta struct {
	ModuleName string `cty:"module_name"`
}

// Config is the configuration structure used to instantiate the Google
// provider.
type Config struct {
	// DCLConfig
	AccessToken                        string
	Credentials                        string
	ImpersonateServiceAccount          string
	ImpersonateServiceAccountDelegates []string
	Project                            string
	Region                             string
	BillingProject                     string
	Zone                               string
	Scopes                             []string
	// BatchingConfig                     *batchingConfig
	UserProjectOverride bool
	RequestReason       string
	RequestTimeout      time.Duration
	// PollInterval is passed to resource.StateChangeConf in common_operation.go
	// It controls the interval at which we poll for successful operations
	PollInterval time.Duration

	client             *http.Client
	context            context.Context
	userAgent          string
	gRPCLoggingOptions []option.ClientOption

	//tokenSource oauth2.TokenSource

	AccessApprovalBasePath       string
	AccessContextManagerBasePath string
	ActiveDirectoryBasePath      string
	ApiGatewayBasePath           string
	ApigeeBasePath               string
	AppEngineBasePath            string
	ArtifactRegistryBasePath     string
	BigQueryBasePath             string
	BigqueryAnalyticsHubBasePath string
	BigqueryConnectionBasePath   string
	BigqueryDatapolicyBasePath   string
	BigqueryDataTransferBasePath string
	BigqueryReservationBasePath  string
	BigtableBasePath             string
	BillingBasePath              string
	BinaryAuthorizationBasePath  string
	CertificateManagerBasePath   string
	CloudAssetBasePath           string
	CloudBuildBasePath           string
	CloudFunctionsBasePath       string
	Cloudfunctions2BasePath      string
	CloudIdentityBasePath        string
	CloudIdsBasePath             string
	CloudIotBasePath             string
	CloudRunBasePath             string
	CloudSchedulerBasePath       string
	CloudTasksBasePath           string
	ComputeBasePath              string
	ContainerAnalysisBasePath    string
	DataCatalogBasePath          string
	DataFusionBasePath           string
	DataLossPreventionBasePath   string
	DataprocBasePath             string
	DataprocMetastoreBasePath    string
	DatastoreBasePath            string
	DatastreamBasePath           string
	DeploymentManagerBasePath    string
	DialogflowBasePath           string
	DialogflowCXBasePath         string
	DNSBasePath                  string
	DocumentAIBasePath           string
	EssentialContactsBasePath    string
	FilestoreBasePath            string
	FirebaseBasePath             string
	FirestoreBasePath            string
	GameServicesBasePath         string
	GKEHubBasePath               string
	HealthcareBasePath           string
	IAM2BasePath                 string
	IAMBetaBasePath              string
	IapBasePath                  string
	IdentityPlatformBasePath     string
	KMSBasePath                  string
	LoggingBasePath              string
	MemcacheBasePath             string
	MLEngineBasePath             string
	MonitoringBasePath           string
	NetworkManagementBasePath    string
	NetworkServicesBasePath      string
	NotebooksBasePath            string
	OrgPolicyBasePath            string
	OSConfigBasePath             string
	OSLoginBasePath              string
	PrivatecaBasePath            string
	PubsubBasePath               string
	PubsubLiteBasePath           string
	RedisBasePath                string
	ResourceManagerBasePath      string
	RuntimeConfigBasePath        string
	SecretManagerBasePath        string
	SecurityCenterBasePath       string
	SecurityScannerBasePath      string
	ServiceDirectoryBasePath     string
	ServiceManagementBasePath    string
	ServiceUsageBasePath         string
	SourceRepoBasePath           string
	SpannerBasePath              string
	SQLBasePath                  string
	StorageBasePath              string
	TagsBasePath                 string
	TPUBasePath                  string
	VertexAIBasePath             string
	VPCAccessBasePath            string
	WorkflowsBasePath            string

	CloudBillingBasePath      string
	ComposerBasePath          string
	ContainerBasePath         string
	DataflowBasePath          string
	IamCredentialsBasePath    string
	ResourceManagerV3BasePath string
	IAMBasePath               string
	CloudIoTBasePath          string
	ServiceNetworkingBasePath string
	StorageTransferBasePath   string
	BigtableAdminBasePath     string

	// dcl
	ContainerAwsBasePath   string
	ContainerAzureBasePath string

	//requestBatcherServiceUsage *RequestBatcher
	//requestBatcherIam          *RequestBatcher
}

// Methods to create new services from config
// Some base paths below need the version and possibly more of the path
// set on them. The client libraries are inconsistent about which values they need;
// while most only want the host URL, some older ones also want the version and some
// of those "projects" as well. You can find out if this is required by looking at
// the basePath value in the client library file.
func (c *Config) NewComputeClient(userAgent string) *compute.Service {
	log.Printf("[INFO] Instantiating GCE client for path %s", c.ComputeBasePath)
	clientCompute, err := compute.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client compute: %s", err)
		return nil
	}
	clientCompute.UserAgent = userAgent
	clientCompute.BasePath = c.ComputeBasePath

	return clientCompute
}

func (c *Config) NewContainerClient(userAgent string) *container.Service {
	containerClientBasePath := removeBasePathVersion(c.ContainerBasePath)
	log.Printf("[INFO] Instantiating GKE client for path %s", containerClientBasePath)
	clientContainer, err := container.NewService(c.context, option.WithHTTPClient(c.client))
	if err != nil {
		log.Printf("[WARN] Error creating client container: %s", err)
		return nil
	}
	clientContainer.UserAgent = userAgent
	clientContainer.BasePath = containerClientBasePath

	return clientContainer
}

// Remove the `/{{version}}/` from a base path if present.
func removeBasePathVersion(url string) string {
	re := regexp.MustCompile(`(?P<base>http[s]://.*)(?P<version>/[^/]+?/$)`)
	return re.ReplaceAllString(url, "$1/")
}
