// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package archive

import (
	"fmt"
	"io"

	"github.com/opencontainers/go-digest"
)

const emptySHA256TarDigest = "sha256:5f70bf18a086007016e948b04aed3b82103a36bea41755b6cddfaf10ace3c6ef"

func TarCanonicalDigest(path string) (digest.Digest, int64, error) {
	tarReader, err := TarWithoutRootDir(path)
	if err != nil {
		return "", 0, fmt.Errorf("unable to tar on %s, err: %s", path, err)
	}
	defer tarReader.Close()

	digester := digest.Canonical.Digester()
	size, err := io.Copy(digester.Hash(), tarReader)
	if err != nil {
		return "", 0, err
	}
	layerDigest := digester.Digest()
	if layerDigest == emptySHA256TarDigest {
		return "", 0, nil
	}

	return layerDigest, size, nil
}
