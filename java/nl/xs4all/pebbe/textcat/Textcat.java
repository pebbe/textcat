package nl.xs4all.pebbe.textcat;

import java.io.File;
import java.io.FileNotFoundException;
import java.util.ArrayList;
import java.util.Collections;
import java.util.Comparator;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.Scanner;

public class Textcat {

    private double thresholdValue;
    private int maxCandidates;
    private int minDocSize;
    private int maxPatterns = 400;
    private String[] langs = new String[60];
    private List<HashMap<String, Integer>> data = new ArrayList<HashMap<String, Integer>>(60);

    private class countType {
        String s;
        int i;
	public countType(String s, int i) {
	    this.s = s;
	    this.i = i;
	}
    }

    public Textcat() {
	thresholdValue = 1.03;
	maxCandidates = 5;
	minDocSize = 25;

	ClassLoader classLoader = getClass().getClassLoader();
	File file = new File(classLoader.getResource("nl/xs4all/pebbe/textcat/data").getFile());
	try {
	    Scanner scanner = new Scanner(file);
	    int state = 0;
	    HashMap<String, Integer> d = new HashMap<String, Integer>();
	    int count = 0;
	    int i = 0;
	    while (scanner.hasNextLine()) {
		String line = scanner.nextLine();
		if (state == 0) {
		    langs[count] = line;
		    i = 0;
		    d = new HashMap<String, Integer>(536, .75f); // 400 items per language
		    state = 1;
		} else if (line.equals("******")) {
		    data.add(d);
		    count++;
		    state = 0;
		} else {
		    d.put(line, i);
		    i++;
		}
	    }
	    scanner.close();
	} catch (FileNotFoundException e) {
	    e.printStackTrace();
	}
    }

    public void setMinDocSize(int i) {
	minDocSize = i;
	if (minDocSize < 1) {
	    minDocSize = 25;
	}
    }

    public void setThresholdValue(double d) {
	thresholdValue = d;
	if (thresholdValue < 1) {
	    thresholdValue = 1.03;
	}
    }

    public void setMaxCandidates(int i) {
	maxCandidates = i;
	if (maxCandidates < 1) {
	    maxCandidates = 5;
	}
    }

    public String[] classify(String text) {

	if (text.length() < minDocSize) {
	    String[] r = new String[1];
	    r[0] = "SHORT";
	    return r;
	}

	List<countType> scores = new ArrayList<countType>(langs.length);

	String[] patt = getPatterns(text);
	for (int idx = 0; idx < langs.length; idx++) {
	    String lang = langs[idx];
	    HashMap<String, Integer> dat = data.get(idx);
	    int score = 0;
	    for (int n = 0; n < patt.length; n++) {
		String p = patt[n];
		int i = maxPatterns;
		if (dat.containsKey(p)) {
		    i = dat.get(p);
		}
		if (n > i) {
		    score += n - i;
		} else {
		    score += i - n;
		}
	    }
	    scores.add(new countType(lang, score));
	}

        int minScore = maxPatterns * maxPatterns;
        for (int i = 0; i < langs.length; i++) {
	    countType ct = scores.get(i);
	    if (ct.i < minScore) {
		minScore = ct.i;
	    }
        }
        double threshold = minScore * thresholdValue;
        int nCandidates = 0;
        for (int i = 0; i < langs.length; i++) {
	    if (scores.get(i).i <= threshold) {
		nCandidates++;
	    }
        }
        if (nCandidates > maxCandidates) {
	    String[] r = new String[1];
	    r[0] = "UNKNOWN";
	    return r;
        }

	List<countType> lowScores = new ArrayList<countType>(nCandidates);
        for (int i = 0; i < langs.length; i++) {
	    countType ct = scores.get(i);
	    if (ct.i <= threshold) {
		lowScores.add(ct);
	    }
        }
	Collections.sort(lowScores, new Comparator<countType>() {
		@Override
		public int compare(countType o1, countType o2) {
		    return o1.i - o2.i;
		}
	    });
	String[] languages = new String[lowScores.size()];
        for (int i = 0; i < languages.length; i++) {
	    languages[i] = lowScores.get(i).s;
        }
	return languages;
    }

    private String[] getPatterns(String s) {
        HashMap<String, Integer> ngrams = new HashMap<String, Integer>(s.length() * 6);
	String[] words = s.toLowerCase().split("[^\\p{L}]+");
	for (int w = 0; w < words.length; w++) {
	    String word = "_" + words[w] + "____";
	    int n = word.length() - 4;
	    for (int i = 0; i < n; i++) {
		for (int j = 1; j < 6; j++) {
		    String ng = word.substring(i, i+j);
		    if (!ng.endsWith("__")) {
			if (!ngrams.containsKey(ng)) {
			    ngrams.put(ng, 0);
			}
			ngrams.put(ng, ngrams.get(ng) + 1);
		    }
		}
	    }
	}
        int size = ngrams.size();
        List<countType> counts = new ArrayList<countType>(size);
	for (Map.Entry<String, Integer> item : ngrams.entrySet()) {
	    counts.add(new countType(item.getKey(), item.getValue()));
	}
	Collections.sort(counts, new Comparator<countType>() {
		@Override
		public int compare(countType o1, countType o2) {
		    return o2.i - o1.i;
		}
	    });
	if (size > maxPatterns) {
	    size = maxPatterns;
	}
	String[] result = new String[size];
	for (int i = 0; i < size; i++) {
	    result[i] = counts.get(i).s;
	}
        return result;
    }
}
