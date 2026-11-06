
package patch

import (
    "strings"
    "testing"
)

// validManifest returns two valid releases for testing.
func validManifest() []Manifest {
    return []Manifest{
        {
            Version: "20.1.23",
            PackageVersion: PackageVersion{
                Baadbaan:  "v20.1.23",
                Technical: "v20.1.16",
                Docker:    "v20.2.1",
                Toolbox:   "v20.1.5",
                Discovery: "v20.1.0",
            },
        },
        {
            Version: "20.1.24",
            PackageVersion: PackageVersion{
                Baadbaan:  "v20.1.24",
                Technical: "v20.1.16",
                Docker:    "v20.2.1",
                Toolbox:   "v20.1.5",
                Discovery: "v20.1.0",
            },
        },
    }
}

func TestValidateManifest(t *testing.T) {
    tests := []struct {
        name          string
        modify        func([]Manifest) []Manifest
        wantError     bool
        errorContains string
    }{
        {
            name: "valid manifest",
            modify: func(m []Manifest) []Manifest {
                return m
            },
            wantError: false,
        },
        {
            name: "empty manifest",
            modify: func(m []Manifest) []Manifest {
                return nil
            },
            wantError:     true,
            errorContains: "expected exactly 2 releases",
        },
        {
            name: "single release",
            modify: func(m []Manifest) []Manifest {
                return m[:1]
            },
            wantError:     true,
            errorContains: "expected exactly 2 releases",
        },
        {
            name: "duplicate versions",
            modify: func(m []Manifest) []Manifest {
                m[1].Version = m[0].Version
                return m
            },
            wantError:     true,
            errorContains: "source and target versions must be different",
        },
        {
            name: "empty release version",
            modify: func(m []Manifest) []Manifest {
                m[0].Version = ""
                return m
            },
            wantError:     true,
            errorContains: "release 1: version is required",
        },
        {
            name: "whitespace release version",
            modify: func(m []Manifest) []Manifest {
                m[1].Version = "   "
                return m
            },
            wantError:     true,
            errorContains: "release 2: version is required",
        },
        {
            name: "empty component version",
            modify: func(m []Manifest) []Manifest {
                m[1].PackageVersion.Technical = ""
                return m
            },
            wantError:     true,
            errorContains: "release 2: package_version.technical is required",
        },
        {
            name: "multiple missing components",
            modify: func(m []Manifest) []Manifest {
                m[0].PackageVersion.Docker = ""
                m[0].PackageVersion.Technical = ""
                m[1].PackageVersion.Discovery = ""
                return m
            },
            wantError:     true,
            errorContains: "release 2: package_version.discovery is required",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            manifests := validManifest()
            manifests = tt.modify(manifests)

            err := ValidateManifest(manifests)

            if tt.wantError {
                if err == nil {
                    t.Fatal("expected an error, got nil")
                }

                if !strings.Contains(err.Error(), tt.errorContains) {
                    t.Fatalf(
                        "expected error containing %q, got %q",
                        tt.errorContains,
                        err.Error(),
                    )
                }

                return
            }

            if err != nil {
                t.Fatalf("expected no error, got: %v", err)
            }
        })
    }
}
